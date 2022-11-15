package main

import (
	//"fmt"
	//"github.com/MuddCreates/hyperschedule-api-go/internal/lingk"
	"github.com/alecthomas/kong"
	"log"
	//"net/http"
)

type Cmd struct {
	Port            uint16 `help:"HTTP port to listen on" default:"8332" env:"PORT"`
	DbUrl           string `help:"URL of PostgreSQL database" env:"DB_URL" required:"true"`
	UploadEmailHash string `help:"SHA256 hash of authorized uploader email" required:"true" env:"UPLOAD_EMAIL_HASH"`
	UploadTokenHash string `help:"SHA256 hash of upload token set in IFTTT webhook headers" env:"UPLOAD_TOKEN_HASH" required:"true"`
}

func (c *Cmd) Run() error {
	srv, err := c.NewServer()
	if err != nil {
		return err
	}

	if err := srv.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	cli := Cmd{}
	kong.Parse(&cli)
	if err := cli.Run(); err != nil {
		log.Fatalf("failed: %v", err)
	}
}
