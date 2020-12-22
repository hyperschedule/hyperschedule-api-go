package main

import (
  //"io"
  "log"
  "io/ioutil"
  "os"
  "github.com/MuddCreates/hyperschedule-api-go/internal/badcsv"
  //"github.com/kr/pretty"
  "github.com/davecgh/go-spew/spew"
)

func main() {
  f, err := os.Open("sample/course_1.csv")
  if err != nil {
    log.Fatalf("couldn't open file: %v", err)
  }

  input, err := ioutil.ReadAll(f)
  if err != nil {
    log.Fatalf("failed to read: %v", err)
  }
  _, warns, fails, err := badcsv.Parse(input)
  if err != nil {
    log.Fatalf("fail: %v", err)
  }
  //pretty.Logln(x)
  //log.Println("nice")
  spew.Dump(warns)
  spew.Dump(fails)


  //_ = f

  /*
	r := fakecsv.New(bufio.NewReader(f))
  for {
	  _, err := r.ReadRow()
    if err == io.EOF {
      break
    }
	  if err != nil {
      log.Fatalf("failed to read csv row: %v", err)
    }
  }
  log.Printf("successfully read all rows")
  */
}
