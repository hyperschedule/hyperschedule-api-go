package coursesectionschedule

import (
  "io"
  "strconv"
  "errors"
  "github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
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

const dayString = "UMTWRFS"

func parseTime(s string) (*Time, error) {
  n, err := strconv.Atoi(s)
  if err != nil {
    return nil, err
  }
  h, m := n / 100, n % 100
  if !(0 <= h && h < 24 && 0 <= m && m < 60) {
    return nil, errors.New("bad time range")
  }
  return &Time{Hour: h, Minute: m}, nil
}

func parseDays(s string) (Days, error) {
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

func parse(record []string) (*CourseSectionSchedule, error) {
  colExternalId := record[0]
  colCourseSectionExternalId := record[1]
  colClassBeginningTime := record[2]
  colClassEndingTime := record[3]
  colClassMeetingDays := record[4]
  colInstructionSiteName := record[5]
  //colInstructionSiteType := record[6]

  start, err := parseTime(colClassBeginningTime)
  if err != nil {
    return nil, err
  }
  end, err := parseTime(colClassEndingTime)
  if err != nil {
    return nil, err
  }
  days, err := parseDays(colClassMeetingDays)
  if err != nil {
    return nil, err
  }

  return &CourseSectionSchedule{
    Id: colExternalId,
    CourseSectionId: colCourseSectionExternalId,
    Start: start,
    End: end,
    Days: days,
    Location: colInstructionSiteName,
  }, nil
}

func ReadAll(r io.Reader) ([]*CourseSectionSchedule, []error, error) {
  courseSectionSchedules := make([]*CourseSectionSchedule, 0, 1024)
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
