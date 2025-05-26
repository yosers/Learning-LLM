create-db:
	docker run --name shofy -e POSTGRES_PASSWORD=papabear -e POSTGRES_DB=shofy -e POSTGRES_USER=postgres -p 5433:5432 -d postgres

drop-db:
	docker rm -f shofy

migrate-up:
	migrate -path db/migration -database "postgresql://postgres:papabear@localhost:5433/shofy?sslmode=disable" up

migrate-down:
	migrate -path db/migration -database "postgresql://postgres:papabear@localhost:5433/shofy?sslmode=disable" down