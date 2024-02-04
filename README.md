## Run API docker image
```
docker run -p 5051:5051 -e BINDER_DB_DRIVER=postgres -e BINDER_DB_SOURCE="postgres://postgres:${DKR_POSTGRES_PWD}@postgres:${DKR_POSTGRES_PORT}/binder_ut?sslmode=disable" -e GIN_MODE=release -e BINDER_API_SERVER_ADDRESS=":5051" -d --name gobinder__api --network gobind-network c0d0df130de5
```

## Run UI docker image
```
docker run -p 5050:5050 -e BINDER_UI_SERVER_ADDRESS=":5050" -e REDIS_URI="redis:6379" -e BINDER_API_BASE_URL="http://gobinder__api:5051" -e BINDER_DB_DRIVER=postgres -d --name gobinder__ui --network gobind-network c1eab1ddd654
```