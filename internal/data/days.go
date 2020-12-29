package data

import (
  "errors"
)

type Days int

const dayString = "UMTWRFS"

func ParseDays(s string) (Days, error) {
  var days Days
  if len(s) != len(dayString) {
    return days, errors.New("bad days length")
  }
  for i, c := range s {
    switch byte(c) {
    case '-':
      case dayString[i]: days |= 1 << i
    default:
      return days, errors.New("unexpected character in day string")
    }
  }
  return days, nil
}
