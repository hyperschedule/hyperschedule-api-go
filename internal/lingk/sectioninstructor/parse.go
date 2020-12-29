package sectioninstructor

import (
  "io"
  "github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
)

var expectHead = []string{
  "courseSectionExternalId",
  "staffExternalId",
}

func parse(record []string) *Entry {
  return &Entry{
    CourseSectionId: record[0],
    StaffId: record[1],
  }
}

func ReadAll(r io.Reader) ([]*Entry, []error, error) {
  entries := make([]*Entry, 0, 2048)
  errs, err := csvutil.Collect(r, expectHead, func(record []string) error {
    entries = append(entries, parse(record))
    return nil
  })
  return entries, errs, err
}
