package main

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

var cmd = &cobra.Command{
	Use:   "hyperschedule-server",
	Short: "API server for Hyperschedule",
	Run:   run,
}
var uploadEmailHash string

func run(cmd *cobra.Command, args []string) {
	log.Printf("loading initial data")
	data, err := lingk.Sample()
	if err != nil {
		log.Fatalf("failed to load: %v", err)
	}
	state.SetData(data)

	addr := fmt.Sprintf(":%s", viper.GetString("port"))
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("bad %v", err)
	}
}

func init() {
	http.HandleFunc("/upload/", inboundHandler)
	http.HandleFunc("/api/v3/", apiV3Handler)
	http.HandleFunc("/raw/", rawHandler)
	http.HandleFunc("/raw/staff/", rawStaffHandler)

	viper.AutomaticEnv()

	cmd.Flags().Int("port", 8332, "HTTP port to listen on.")
	viper.BindPFlag("port", cmd.Flags().Lookup("port"))

	uploadEmailHash = os.Getenv("UPLOAD_EMAIL_HASH")
	if len(uploadEmailHash) == 0 {
		log.Printf("warning: did not define UPLOAD_EMAIL_HASH")
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
