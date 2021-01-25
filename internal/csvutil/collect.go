package csvutil

import (
  "io"
  "encoding/csv"
  "errors"
  "fmt"
)

func Collect(
  r io.Reader,
  expectHead []string,
  f func([]string) error,
) ([]error, error) {
  reader := csv.NewReader(r)

  head, err := reader.Read()
  if err != nil {
    return nil, errors.New("failed to read header")
  }
  if len(head) != len(expectHead) {
    return nil, errors.New("header length mismatch")
  }
  for i, h := range head {
    if h != expectHead[i] {
      return nil, errors.New(fmt.Sprintf("mismatch header: expecting %s but got %s on column %d", expectHead[i], h, i))
    }
  }

  errs := make([]error, 0)
  for {
    record, err := reader.Read()
    if err != nil {
    if err == io.EOF {
      break
    }
      return nil, err
    }
    err = f(record)
    if err != nil {
      errs = append(errs, err)
    }
  }

  return errs, nil
}
