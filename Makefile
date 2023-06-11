
postgres:
	docker run --rm \
    --detach \
    --publish 5432:5432 \
    --env POSTGRES_DB=postgres \
    --env POSTGRES_USER=postgres \
    --env POSTGRES_PASSWORD=postgres \
    postgres

migrate:
	migrate -path ./migrations/postgres/ -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up