package data

import (
  "strconv"
  "errors"
  "fmt"
)

type Time struct{
  Hour int
  Minute int
}

func ParseTime(s string) (Time, error) {
  n, err := strconv.Atoi(s)
  if err != nil {
    return Time{}, err
  }
  h, m := n / 100, n % 100
  if !(0 <= h && h < 24 && 0 <= m && m < 60) {
    return Time{}, errors.New("bad time range")
  }
  return Time{Hour: h, Minute: m}, nil
}

func (t Time) String() string {
  return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}
