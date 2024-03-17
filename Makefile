DB_URL=postgres://root:IWSIWDF2024@localhost:5432/legal_referral_db?sslmode=disable

postgres:
	docker run --name legal-referral-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=IWSIWDF2024 -d postgres:16.2-alpine3.19

createdb:
	docker exec -it legal-referral createdb --username=root --owner=root legal_referral_db

dropdb:
	docker exec -it legal-referral dropdb legal_referral_db

sqlc:
	sqlc generate

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

server:
	go run main.go

.PHONY: postgres createdb dropdb sqlc migrateup migratedown
