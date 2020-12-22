package csvutil

import (
  "io"
  "encoding/csv"
  "errors"
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
      return nil, errors.New("mismatch header")
    }
  }

  errs := make([]error, 0)
  for {
    record, err := reader.Read()
    if err != io.EOF {
      break
    }
    if err != nil {
      return nil, err
    }
    err = f(record)
    if err != nil {
      errs = append(errs, err)
    }
  }

  return errs, nil
}
