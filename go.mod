module github.com/MuddCreates/hyperschedule-api-go

// +heroku goVersion go1.16
// +heroku install -tags 'postgres' ./vendor/github.com/golang-migrate/migrate/v4/cmd/migrate .
go 1.16

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.6.0 // indirect
	git.apache.org/thrift.git v0.0.0-20180924222215-a9235805469b // indirect
	github.com/cznic/ql v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.39.0 // indirect
	github.com/golang-migrate/migrate/v4 v4.14.1 // indirect
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/google/pprof v0.0.0-20201109224723-20978b51388d // indirect
	github.com/googleapis/gax-go v2.0.0+incompatible // indirect
	github.com/gotestyourself/gotestyourself v2.1.0+incompatible // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.2.0+incompatible // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kshvakov/clickhouse v1.3.4 // indirect
	github.com/openzipkin/zipkin-go v0.1.1 // indirect
	github.com/sendgrid/rest v2.6.2+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.7.2+incompatible // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	golang.org/x/arch v0.0.0-20201008161808-52c3e6f60cff // indirect
	golang.org/x/crypto v0.0.0-20201112155050-0c6587e931a9 // indirect
	golang.org/x/sys v0.0.0-20201113233024-12cec1faf1ba // indirect
)
