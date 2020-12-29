package calendarsession

type Date struct{
  Year int
  Month int
  Day int
}

type Entry struct{
  Id string
  Semester string
  Start Date
  End Date
}
