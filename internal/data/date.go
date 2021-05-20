package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

var defaultMonthDays = []int{
	31, // jan
	28, // feb
	31, // mar
	30, // apr
	31, // may
	30, // jun
	31, // jul
	31, // aug
	30, // sep
	31, // oct
	30, // nov
	31, // dec
}

func leap(y, m int) int {
	if m == 2 && (y%4 == 0 && y%100 != 0 || y%400 == 0) {
		return 1
	} else {
		return 0
	}
}

func monthDays(y, m int) int {
	return defaultMonthDays[m-1] + leap(y, m)
}

func ParseDate(s string) (Date, error) {
	segs := strings.Split(s, "-")
	if len(segs) != 3 {
		return Date{}, errors.New("wrong number of segments in date")
	}
	year, err := strconv.Atoi(segs[0])
	if err != nil {
		return Date{}, errors.New("bad year int")
	}
	month, err := strconv.Atoi(segs[1])
	if err != nil {
		return Date{}, errors.New("bad month int")
	}
	day, err := strconv.Atoi(segs[2])
	if err != nil {
		return Date{}, errors.New("bad day int")
	}

	if !(2000 <= year && year < 2050) {

		// *Technically*, this isn't very future proof, but 30 years from now is a
		// long way to go---if any of our entries appear to contain a year greater
		// than 2050, chances are that it's not actually the future; it's just a
		// bug in the formatting we're expecting.
		return Date{}, errors.New("bad year range")
	}

	if !(1 <= month && month <= 12) {
		return Date{}, errors.New("bad month")
	}

	if !(1 <= day && day <= monthDays(year, month)) {
		return Date{}, errors.New("bad days")
	}

	return Date{
		Year:  year,
		Month: month,
		Day:   day,
	}, nil
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d Date) ToTime() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}
