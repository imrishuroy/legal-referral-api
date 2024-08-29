#DB_URL=postgres://root:IWSIWDF2024@localhost:5432/legal_referral_db?sslmode=disable

#DB_URL=postgres://root:IWSIWDF2024@localhost:5432/legal_referral_db?sslmode=disable

DB_URL=postgres://root:J01J7QiaHr3apJgUYpT6@legal-referral.ct2smiqa0pnv.us-east-1.rds.amazonaws.com/legal_referral

postgres:
	docker run --name legal-referral-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=IWSIWDF2024 -d postgres:16.4-alpine3.20

redis:
	docker run --name legal-referral-redis -p 6379:6379 -d redis:7-alpine

createdb:
	docker exec -it legal-referral-db createdb --username=root --owner=root legal_referral_db

dropdb:
	docker exec -it legal-referral dropdb legal_referral_db

sqlc:
	sqlc generate

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migratedown2:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 8

server:
	go run main.go

.PHONY: postgres createdb dropdb sqlc migrateUp migrateDown migrateDown2 new_migration
