package main

import (
  "os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var cmd = &cobra.Command{
	Use:   "hyperschedule-server",
	Short: "API server for Hyperschedule",
	Run: run,
}
var uploadEmailHash string

func run(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf(":%s", viper.GetString("port"))
  log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("bad %v", err)
	}
}

func init() {
	http.HandleFunc("/upload/", inboundHandler)
  http.HandleFunc("/api/v3/", apiV3Handler)

  viper.AutomaticEnv()
	cmd.Flags().Int("port", 8332, "HTTP port to listen on.")
  viper.BindPFlag("port", cmd.Flags().Lookup("port"))

  uploadEmailHash = os.Getenv("UPLOAD_EMAIL_HASH")
  if len(uploadEmailHash) == 0 {
    log.Fatalf("forgot to define UPLOAD_EMAIL_HASH")
  }
}

func main() {
  //if err := viper.ReadInConfig(); err != nil {
  //  log.Fatalf("failed to read config, %v", err)
  //}
	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed %v", err)
	}
}
