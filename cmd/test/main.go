package main

import (
  "log"
  "bufio"
  "os"
  "github.com/MuddCreates/hyperschedule-api-go/internal/fakecsv"
)

func main() {
	r := fakecsv.New(bufio.NewReader(os.Stdin))
	row, err := r.ReadRow()
	if err != nil { log.Fatalf("aaa %v", err) }
	log.Printf("%#v\n", row)
	//log.Println(parseFakeCsv(r))
}
