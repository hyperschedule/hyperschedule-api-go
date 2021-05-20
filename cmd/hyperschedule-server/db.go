package main

import (
	"context"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"log"
	"os"
)

var dbConn *db.Connection

//func dbConnect(addr string) error {
//	//db.New(context.Background(), )
//
//}

func init() {
	url := os.Getenv("DB_URL")
	if len(url) == 0 {
		log.Fatalf("missing db url")
	}

	conn, err := db.New(context.Background(), url)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	dbConn = conn
}
