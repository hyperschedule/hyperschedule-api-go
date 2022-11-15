package data

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Time struct {
	Hour   int
	Minute int
}

func ParseTime(s string) (Time, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return Time{}, err
	}
	h, m := n/100, n%100
	if !(0 <= h && h < 24 && 0 <= m && m < 60) {
		return Time{}, errors.New("bad time range")
	}
	if h == 0 {
		h = 12
	}
	return Time{Hour: h, Minute: m}, nil
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

func (t Time) Std() time.Time {
	return time.Date(0, time.January, 0, t.Hour, t.Minute, 0, 0, time.UTC)
}

func TimeFromStd(t time.Time) Time {
	return Time{t.Hour(), t.Minute()}
}

func TimeLess(t1, t2 Time) bool {
	return t1.Hour < t2.Hour || t1.Hour == t2.Hour && t1.Minute < t2.Minute
}
