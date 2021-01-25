package main

import (
	"fmt"
	//"github.com/MuddCreates/hyperschedule-api-go/internal/lingk"
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

func run(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("bad %v", err)
	}
}

func init() {
	http.HandleFunc("/upload/", inboundHandler)
	cmd.Flags().IntVar(&port, "port", 80, "HTTP port to listen on.")
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed %v", err)
	}
}
