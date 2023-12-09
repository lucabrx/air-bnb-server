db/container/create:
	docker run --name air-bnb -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres

db/migration/create:
	echo "Creating migration... $(name)"
	migrate create -seq -ext=.sql -dir=./migrations "$(name)"

db/migration/up:
	echo "Migrating up..."
	migrate -path ./migrations -database "$(DB_URL)"  up

db/migration/ci/test:
	echo "Migrating up..."
	migrate -path ./migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable  up

db/services/test:
	go test -v ./...


.PHONY: db/container/create db/migration/create db/migration/up db/migration/ci/test db/services/test