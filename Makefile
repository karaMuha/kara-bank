appup:
	docker-compose up -d

migrateup:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose down

migratedown1:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose down 1

appdown:
	docker-compose down

testall:
	cd cmd; \
	go test -v ./...