DB_DOCKER_CONTAINER=friend-management-db-1
BINARY_NAME=app
build:
	docker-compose build

run:
	docker-compose up

stop:
	docker-compose down

run_migrations:
	docker cp data/migrations/. ${DB_DOCKER_CONTAINER}:/migrations
	docker exec -it ${DB_DOCKER_CONTAINER} psql -U friend-management -d friend-management -f /migrations/001_setup_db.down.sql
	docker exec -it ${DB_DOCKER_CONTAINER} psql -U friend-management -d friend-management -f /migrations/001_setup_db.up.sql 

create_migrations:
	sqlboiler psql -c sqlboiler.toml --wipe --no-tests

test:
	go test -mod=vendor -coverprofile=c.out -failfast -timeout 5m ./...

