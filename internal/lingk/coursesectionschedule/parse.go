package coursesectionschedule

import (
  "io"
  "strings"
  "github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
  "github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

var expectHead = []string{
  "externalId",
  "courseSectionExternalId",
  "classBeginningTime",
  "classEndingTime",
  "classMeetingDays",
  "InstructionSiteName",
  "InstructionSiteType",
}


func parse(record []string) (*Entry, error) {
  colExternalId := record[0]
  colCourseSectionExternalId := record[1]
  colClassBeginningTime := record[2]
  colClassEndingTime := record[3]
  colClassMeetingDays := record[4]
  colInstructionSiteName := record[5]
  //colInstructionSiteType := record[6]

  start, err := data.ParseTime(colClassBeginningTime)
  if err != nil {
    return nil, err
  }
  end, err := data.ParseTime(colClassEndingTime)
  if err != nil {
    return nil, err
  }
  days, err := data.ParseDays(colClassMeetingDays)
  if err != nil {
    return nil, err
  }

  return &Entry{
    Id: colExternalId,
    CourseSectionId: colCourseSectionExternalId,
    Start: start,
    End: end,
    Days: days,
    Location: strings.TrimSpace(colInstructionSiteName),
  }, nil
}

func ReadAll(r io.Reader) ([]*Entry, []error, error) {
  courseSectionSchedules := make([]*Entry, 0, 1024)
  errs, err := csvutil.Collect(r, expectHead, func(record []string) error {
    courseSectionSchedule, err := parse(record)
    if err != nil {
      return err
    }
    courseSectionSchedules = append(courseSectionSchedules, courseSectionSchedule)
    return nil
  })
  return courseSectionSchedules, errs, err
}
