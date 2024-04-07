DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

# docker
postgres:
	 docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e TZ=Asia/Shanghai -d postgres:alpine3.18
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres16 dropdb simple_bank
stop:
	docker stop postgres16
start:
	docker start postgres16
# migrate up or down
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

# sqlc
sqlc:
	sqlc generate

# go
test:
	go test -v -cover -short ./...

.PHONY: createdb dropdb postgres migratedown migrateup sqlc stop start test
