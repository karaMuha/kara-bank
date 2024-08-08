appup:
	docker-compose up -d

migrateup:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose down

appdown:
	docker-compose down