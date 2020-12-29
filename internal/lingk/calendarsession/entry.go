package calendarsession

import "github.com/MuddCreates/hyperschedule-api-go/internal/data"

type Entry struct{
  Id string
  Semester string
  Start data.Date
  End data.Date
}
