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

migrateup1:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" -verbose up 1

migrateuput:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" -verbose up

migrateuput1:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" -verbose down 1

migratedownut:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" -verbose down

migratedownut1:
	migrate -path db/migrations -database "postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" -verbose down

sqlc:
	sqlc generate

ut:
	BINDER_DB_DRIVER=postgres BINDER_DB_SOURCE="postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" go test -timeout 60s -cover ./...

serve:
	BINDER_DB_DRIVER=postgres BINDER_DB_SOURCE="postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" go run main.go

mock:
	mockgen -package mockdb -destination db/mock/storage.go github.com/dpomian/gobind/db/sqlc Storage

.PHONY: postgres createdb createdbut migrateup migratedown migrateup1 migratedown1 migrateuput migratedownut migrateuput1 migratedownut1 sqlc ut serve mock