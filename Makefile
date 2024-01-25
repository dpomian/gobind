postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=${DKR_POSTGRES_PWD} -v /Users/dpomian/storage/postresql:/var/lib/postgresql/data -p 5444:5432 -d postgres:16.1

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres binder

dropdb:
	docker exec -it postgres dropdb --username=postgres binder

migrateup:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:5444/binder?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:5444/binder?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb migrateup migratedow sqlc