SHELL := /bin/bash -o pipefail

.PHONY: dev
dev:
	set -a && . .env && go run ./cmd/hyperschedule-server
