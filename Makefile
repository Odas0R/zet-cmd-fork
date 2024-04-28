build:
	go build -tags "fts5" -o zet ./cmd
new:
	@read -p "Enter the name of the new migration: " name; \
	go run -tags "fts5" ./cmd/main.go migrate create $$name
up:
	@go run -tags "fts5" ./cmd/main.go migrate up
down:
	@go run -tags "fts5" ./cmd/main.go migrate down
status:
	@go run -tags "fts5" ./cmd/main.go migrate status
schema:
	sqlite3 ./zettel.db .schema
test:
	go test -v -tags "fts5" ./...
