package lingk

import (
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/course"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesection"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesectionschedule"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsession"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsessionsection"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/sectioninstructor"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/staff"
)

type tables struct{
  course []*course.Entry
  courseSection []*coursesection.Entry
  courseSectionSchedule []*coursesectionschedule.Entry
  calendarSession []*calendarsession.Entry
  calendarSessionSection []*calendarsessionsection.Entry
  sectionInstructor []*sectioninstructor.Entry
  staff []*staff.Entry
}

type Entry struct{
}

func (t *tables) join() {
  courses := make(map[string]*course.Entry)
  for _, course := range t.course {
    // TODO check no duplicates
    courses[course.Id] = course
  }

  courseSections := make(map[string]*coursesection.Entry)
  for _, courseSection := range t.courseSections {
    courseSections[courseSection.Id] = courseSection
  }
}
