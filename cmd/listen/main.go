package main

import (
  "fmt"
	//"github.com/MuddCreates/hyperschedule-api-go/internal/lingk"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func initHttp() {
	http.HandleFunc("/upload/", inboundHandler)

}

func initCli() *cobra.Command {
	var port *int

	cmd := &cobra.Command{
		Use:   "hyperschedule-server",
		Short: "API server for Hyperschedule",
		Run: func(cmd *cobra.Command, args []string) {
      addr := fmt.Sprintf(":%d", *port)
			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Fatalf("bad %v", err)
			}
		},
	}

	port = cmd.Flags().Int("port", 80, "HTTP port to listen on.")
	return cmd
}

func main() {
  initHttp()
  if err := initCli().Execute(); err != nil {
    log.Fatalf("failed %v", err)
  }
}
