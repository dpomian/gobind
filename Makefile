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

apiserve:
	BINDER_API_SERVER_ADDRESS=":5056" BINDER_DB_DRIVER=postgres BINDER_DB_SOURCE="postgres://postgres:${DKR_POSTGRES_PWD}@localhost:${DKR_POSTGRES_PORT}/binder?sslmode=disable" go run main_api.go

uiserve:
	BINDER_UI_SERVER_ADDRESS=":5055" BINDER_API_BASE_URL="http://localhost:5056" REDIS_URI=localhost:6379 go run ui/main_ui.go

mock:
	mockgen -package mockdb -destination db/mock/storage.go github.com/dpomian/gobind/db/sqlc Storage

dockerapi:
	docker build -f Dockerfile_api -t gobind__api:latest .

dockerui:
	docker build -f Dockerfile_ui -t gobind__ui:latest .

.PHONY: postgres createdb createdbut migrateup migratedown migrateup1 migratedown1 migrateuput migratedownut migrateuput1 migratedownut1 sqlc ut serve mock uiserve dockerapi dockerui