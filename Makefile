SHELL := /bin/bash -o pipefail

.PHONY: dev
dev:
	set -a && . .env && go run ./cmd/hyperschedule-server

.PHONY: dev-pgcli
dev-pgcli:
	pgcli "$$(heroku config:get 'DEV_URL')"

.PHONY: dev-migrate-up
dev-migrate-up:
	migrate -source 'file://migrate' -database "$$(heroku config:get 'DEV_URL')" up

.PHONY: dev-migrate-down
dev-migrate-down:
	migrate -source 'file://migrate' -database "$$(heroku config:get 'DEV_URL')" down 1

.PHONY: dev-migrate-drop
dev-migrate-drop:
	migrate -source 'file://migrate' -database "$$(heroku config:get 'DEV_URL')" drop

.PHONY: pgcli-prod
pgcli-prod:
	pgcli "$$(heroku config:get 'DATABASE_URL')"

.PHONY: migrate-create
migrate-create:
	migrate 
