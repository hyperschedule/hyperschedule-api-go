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
	case "C":
		return Closed, nil
	case "O":
		return Open, nil
	case "R":
		return Reopened, nil
	default:
		return 0, errors.New("invalid status")
	}
}

func (s Status) String() string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case Reopened:
		return "reopened"
	default:
		return ""
	}
}

func StatusFromString(s string) (Status, error) {
	switch s {
	case "closed":
		return Closed, nil
	case "open":
		return Open, nil
	case "reopened":
		return Reopened, nil
	default:
		return 0, errors.New("invalid status")
	}
}
