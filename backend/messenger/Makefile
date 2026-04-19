db-migrate:
	mkdir -p migrations
	atlas migrate diff dev --env dev --to "file://migrations" --dir "messanger_v001.sql"

db-push:
	atlas schema apply --env dev --url "docker://postgres/15?search_path=messanger"

proto db:
	make db-migrate
	make db-push
