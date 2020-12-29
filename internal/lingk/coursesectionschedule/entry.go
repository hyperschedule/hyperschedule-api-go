package coursesectionschedule

import "github.com/MuddCreates/hyperschedule-api-go/internal/data"

type Entry struct{
  Id string
  CourseSectionId string
  Start data.Time
  End data.Time
  Days data.Days
  Location string
}
