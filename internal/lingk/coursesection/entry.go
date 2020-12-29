package coursesection

import "github.com/MuddCreates/hyperschedule-api-go/internal/data"

type Entry struct{
  Id string
  CourseId string
  Section int
  Seats data.Seats
  Status data.Status
  QuarterCredits int
}
