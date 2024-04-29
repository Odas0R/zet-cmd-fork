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
