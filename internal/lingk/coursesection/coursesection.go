package coursesection


type Status int
const (
  Closed Status = iota
  Open
  Reopened
)

type CourseSection struct{
  Id string
  CourseId string
  Section int
  SeatCapacity int
  SeatEnrolled int
  Status Status
  QuarterCredits int
}
