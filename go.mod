module github.com/MuddCreates/hyperschedule-api-go

go 1.16

// +heroku goVersion go1.16
// +heroku install -tags 'postgres' ./vendor/github.com/golang-migrate/migrate/v4/cmd/migrate .

require (
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/kr/pretty v0.2.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	golang.org/x/sys v0.0.0-20201113233024-12cec1faf1ba // indirect
)
