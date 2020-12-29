package lingk

type Data struct{
  courseSections map[string]*CourseSection
  courses map[string]*Course
  terms map[string]*Term
  staff map[string]Name
}

type CourseSection struct{
  Course string
  Term string
  Section int
  Seats Seats
  Status Status
  QuarterCredits int
  Schedule []*Schedule
  Staff []string
}

type Course struct{
  Title string
  Campus string
  Description string
}

type Term struct{
  Semester string
  Start Date
  End Date
}

type Name struct{
  First string
  Last string
}

type Schedule struct{
  Days Days
  Start Time
  End Time
  Location string
}

type Seats struct{
  Capacity int
  Enrolled int
}

type Days int

type Time struct{
  Hour int
  Minute int
}

type Date struct{
  Year int
  Month int
  Day int
}

type Status int

const (
  Closed Status = iota
  Open
  Reopened
)
