set positional-arguments

dev_postgres_password := 'cool orange banana peels'
dev_upload_email := 'dev@api.hyperschedule.io'

pass:

setup:
  podman pod rm -if hyperschedule-dev
  podman play kube dev-pod.yml

dev *args:
  podman pod start hyperschedule-dev
  go run ./cmd/hyperschedule-server {{args}}

pgcli:
  pgcli "$DB_URL"

migrate *args:
  migrate -path 'migrate' -database "$DB_URL?sslmode=disable" "$@"

upload path:
  zip -qj - '{{path}}'/* \
  | curl \
    -F 'envelope={"from": "", "to": ["{{dev_upload_email}}"]}' \
    -F 'x=@-;filename=HMCarchive.zip' \
    'localhost:8332/upload/'
