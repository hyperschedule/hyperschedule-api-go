package data

type Data struct{
  CourseSections map[string]*CourseSection
  Courses map[string]*Course
  Terms map[string]*Term
  Staff map[string]Name
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

