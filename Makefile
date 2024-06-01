build:
	go build -o zet ./cmd
new:
	@read -p "Enter the name of the new migration: " name; \
	go run ./cmd/main.go migrate create $$name
up:
	@go run ./cmd/main.go migrate up
down:
	@go run ./cmd/main.go migrate down
status:
	@go run ./cmd/main.go migrate status
schema:
	sqlite3 ./zettel.db .schema
serve:
	go run ./cmd/main.go serve
generate:
	templ generate
test:
	go test -v ./...

# ###############################
# Live reload
# ###############################

# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:3777" --proxyport="3000" --open-browser=false

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/cosmtrek/air@v1.52.0 \
	--build.cmd "go build -o tmp/bin/main ./cmd/main.go" --build.bin "tmp/bin/main serve --dev --port 3777" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live:
	make -j2 live/templ live/server
