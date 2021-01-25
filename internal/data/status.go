package data

import "errors"

type Status int

const (
  Closed Status = iota
  Open
  Reopened
)

func ParseStatus(s string) (Status, error) {
  switch s {
  case "C": return Closed, nil
  case "O": return Open, nil
  case "R": return Reopened, nil
  default:
    return 0, errors.New("invalid status")
  }
}

func (s Status) String() string {
  switch s {
  case Closed: return "Closed"
  case Open: return "Open"
  case Reopened: return "Reopened"
  default: return ""
  }
}
