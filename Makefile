db/container/create:
	docker run --name web-dev-tools -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres

db/migration/create:
	echo "Creating migration... $(name)"
	migrate create -seq -ext=.sql -dir=./migrations "$(name)"

db/migration/up:
	echo "Migrating up..."
	migrate -path ./migrations -database "$(DB_URL)"  up

