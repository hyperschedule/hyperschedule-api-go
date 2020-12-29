package calendarsession

import (
  "io"
  "github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
  "strconv"
  "strings"
  "errors"
)

var expectHead = []string{
  "externalId",
  "designator",
  "beginDate",
  "endDate",
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
  if (m == 2 && (y % 4 == 0 && y % 100 != 0 || y % 400 == 0)) {
    return 1
  } else {
    return 0
  }
}

func monthDays(y, m int) int {
  return defaultMonthDays[m-1] + leap(y, m)
}

func parseDate(s string) (*Date, error) {
  segs := strings.Split(s, "-")
  if len(segs) != 3 {
    return nil, errors.New("wrong number of segments in date")
  }
  year, err := strconv.Atoi(segs[0])
  if err != nil {
    return nil, errors.New("bad year int")
  }
  month, err := strconv.Atoi(segs[1])
  if err != nil {
    return nil, errors.New("bad month int")
  }
  day, err := strconv.Atoi(segs[2])
  if err != nil {
    return nil, errors.New("bad day int")
  }

  if !(2000 <= year && year < 2050) {

    // *Technically*, this isn't very future proof, but 30 years from now is a
    // long way to go---if any of our entries appear to contain a year greater
    // than 2050, chances are that it's not actually the future; it's just a
    // bug in the formatting we're expecting.
    return nil, errors.New("bad year range")
  }

  if !(1 <= month && month <= 12) {
    return nil, errors.New("bad month")
  }

  if !(1 <= day && day <= monthDays(year, month)) {
    return nil, errors.New("bad days")
  }

  return &Date{
    Year: year,
    Month: month,
    Day: day,
  }, nil
}

func parse(record []string) (*Entry, error) {
  colExternalId := record[0]
  colDesignator := record[1]
  colBeginDate := record[2]
  colEndDate := record[3]

  start, err := parseDate(colBeginDate)
  if err != nil {
    return nil, err
  }
  end, err := parseDate(colEndDate)
  if err != nil {
    return nil, err
  }

  return &Entry{
    Id: colExternalId,
    Semester: colDesignator,
    Start: start,
    End: end,
  }, nil
}

func ReadAll(r io.Reader) ([]*Entry, []error, error) {
  entries := make([]*Entry, 0, 8)
  errs, err := csvutil.Collect(r, expectHead, func(record []string) error {
    entry, err := parse(record)
    if err != nil {
      return err
    }
    entries = append(entries, entry)
    return nil
  })
  return entries, errs, err
}
