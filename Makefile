postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:Simplepassword@simpo.czhja9t9lf7f.us-east-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:MD8odSFRYfyHB2c2lEM7@simpo.czhja9t9lf7f.us-east-1.rds.amazonaws.com:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:MD8odSFRYfyHB2c2lEM7@simpo.czhja9t9lf7f.us-east-1.rds.amazonaws.com:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:MD8odSFRYfyHB2c2lEM7@simpo.czhja9t9lf7f.us-east-1.rds.amazonaws.com:5432/simple_bank?sslmode=disable" -verbose down 1

migrateforce:
	migrate -path db/migration -database "postgresql://root:MD8odSFRYfyHB2c2lEM7@simpo.czhja9t9lf7f.us-east-1.rds.amazonaws.com:5432/simple_bank?sslmode=disable" -verbose force 1

dropdb:
	docker exec -it postgres12 dropdb simple_bank

sqlc:
	sqlc generate

test:
	go test -v ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migratedown migrateup migrateup1 migratedown1 sqlc test server mock
