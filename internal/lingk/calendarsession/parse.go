package calendarsession

import (
  "io"
  "github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
  "github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

var expectHead = []string{
  "externalId",
  "designator",
  "beginDate",
  "EndDate",
}

func parse(record []string) (*Entry, error) {
  colExternalId := record[0]
  colDesignator := record[1]
  colBeginDate := record[2]
  colEndDate := record[3]

  start, err := data.ParseDate(colBeginDate)
  if err != nil {
    return nil, err
  }
  end, err := data.ParseDate(colEndDate)
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
