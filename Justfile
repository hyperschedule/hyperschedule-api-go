set positional-arguments

dev_postgres_password := 'cool orange banana peels'
dev_upload_email := 'dev@api.hyperschedule.io'

pass:

setup:
  podman pod rm -if hyperschedule-dev
  podman play kube dev-pod.yml

pod:
  podman pod start hyperschedule-dev

dev *args:
  go run ./cmd/hyperschedule-server {{args}}

pgcli:
  pgcli "$DB_URL"

migrate *args:
  migrate -path 'migrate' -database "$DB_URL?sslmode=disable" "$@"

migrate-new name:
  migrate create -ext 'sql' -dir 'migrate' "$1"

upload path:
  @zip -qj - '{{path}}'/* \
  | curl -f \
    -F 'envelope={"from": "", "to": ["{{dev_upload_email}}"]}' \
    -F 'x=@-;filename=HMCarchive.zip' \
    'localhost:8332/upload/'

prod-migrate *args:
  migrate -path 'migrate' -database "$(heroku config:get 'DB_URL')" "$@"

prod-pgcli:
  pgcli "$(heroku config:get 'DB_URL')"
