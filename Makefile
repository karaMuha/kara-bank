init:
	mkdir db-data; \
	mkdir db-data/postgres

buf-update:
	cd cmd; \
	buf dep update

generate-pb:
	cd cmd; \
	buf generate

start:
	docker-compose up -d

stop:
	docker-compose down

migrateup:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose down

migratedown1:
	migrate -path cmd/db/migration -database "postgresql://postgres:postgres@localhost:5433/kara_bank_db?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir ./cmd/db/migration -seq ${name}

testall:
	cd cmd; \
	go test -v ./controllers; \
	go test -v ./db/repositories

proto:
	rm -f cmd/pb/*.go; \
	protoc --proto_path=cmd/proto --go_out=cmd/pb --go_opt=paths=source_relative \
	--go-grpc_out=cmd/pb --go-grpc_opt=paths=source_relative cmd/proto/*.proto