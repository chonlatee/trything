gen: gen-sqlc gen-db

run:
	go run cmd/api/main.go


gen-db:
	go run cmd/gen/*.go

gen-sqlc:
	sqlc generate
