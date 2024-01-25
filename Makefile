postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=${DKR_POSTGRES_PWD} -v /Users/dpomian/storage/postresql:/var/lib/postgresql/data -p 5444:5432 -d postgres:16.1

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres binder

createdbut:
	docker exec -it postgres createdb --username=postgres --owner=postgres binder_ut

dropdb:
	docker exec -it postgres dropdb --username=postgres binder

dropdbut:
	docker exec -it postgres dropdb --username=postgres binder_ut

migrateup:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" -verbose up

migrateuput:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" -verbose down

migratedownut:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" -verbose down

sqlc:
	sqlc generate

ut:
	BINDER_DB_DRIVER=postgres BINDER_DB_SOURCE="postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" go test -timeout 30s -cover github.com/dpomian/gobind/db/sqlc

serve:
	BINDER_DB_DRIVER=postgres BINDER_DB_SOURCE="postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" go run main.go

.PHONY: postgres createdb migrateup migratedow sqlc ut serve