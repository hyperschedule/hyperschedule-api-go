package main

import (
  "os"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var cmd = &cobra.Command{
	Use:   "hyperschedule-server",
	Short: "API server for Hyperschedule",
	Run: run,
}
var port int
var uploadEmailHash string

func run(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("bad %v", err)
	}
}

func init() {
	http.HandleFunc("/upload/", inboundHandler)
  http.HandleFunc("/api/v3/", apiV3Handler)
	cmd.Flags().IntVar(&port, "port", 80, "HTTP port to listen on.")

  uploadEmailHash = os.Getenv("UPLOAD_EMAIL_HASH")
  if len(uploadEmailHash) == 0 {
    log.Fatalf("forgot to define UPLOAD_EMAIL_HASH")
  }
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed %v", err)
	}
}
