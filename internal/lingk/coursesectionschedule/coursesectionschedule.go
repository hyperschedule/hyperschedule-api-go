package coursesectionschedule

type Time struct{
  Hour int
  Minute int
}

type Days int

type CourseSectionSchedule struct{
  Id string
  CourseSectionId string
  Start *Time
  End *Time
  Days Days
  Location string
}
