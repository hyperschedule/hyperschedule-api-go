package lingk

import (
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/course"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesection"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesectionschedule"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsession"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsessionsection"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/sectioninstructor"
  "github.com/MuddCreates/hyperschedule-api-go/internal/lingk/staff"
  "github.com/MuddCreates/hyperschedule-api-go/internal/data"
  "errors"
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

func (t *tables) prune() (*data.Data, []error) {
  p := &data.Data{
    Terms: make(map[string]*data.Term),
    Courses: make(map[string]*data.Course),
    CourseSections: make(map[string]*data.CourseSection),
    Staff: make(map[string]data.Name),
  }
  errs := make([]error, 0)

  // We don't add directly to `p.Courses` here because the raw `course_1.csv`
  // table contains all sorts of extraneous (bad) entries that are never
  // actually referenced by the `coursesection_1.csv` entries--which are the
  // items we actually care about, and which are far more
  // well-formed/consistent than the `course_1.csv` entries.  

  // So we build an intermediate dictionary `courses` below, but only add them
  // to our final dictionary `p.Courses` when they _do_ get referenced by a
  // `coursesection` entry.

  // The `calendarsession_1.csv` and `staff_1.csv` tables do not suffer the
  // bad-data issues as `course_1.csv` does, but we might as well do the same
  // lazy-loading thing, so that we can exclude extraneous entries, just in
  // case there are any.

  courses := make(map[string]*course.Entry, 512)
  for _, c := range t.course {
    if _, ok := courses[c.Id]; ok {
      errs = append(errs, errors.New("course with duplicate id"))
    }
    courses[c.Id] = c
  }

  terms := make(map[string]*calendarsession.Entry, 8)
  for _, c := range t.calendarSession {
    if _, ok := terms[c.Id]; ok {
      errs = append(errs, errors.New("term duplicate id"))
    }
    terms[c.Id] = c
  }

  staff := make(map[string]*staff.Entry, 1024)
  for _, s := range t.staff {
    if _, ok := staff[s.Id]; ok {
      errs = append(errs, errors.New("staff duplicate id"))
    }
    staff[s.Id] = s
  }

  csTerms := make(map[string]string)
  for _, c := range t.calendarSessionSection {
    if _, ok := terms[c.Id]; !ok {
      errs = append(errs, errors.New("calendarsessionsection points to nonexistent term"))
      continue
    }
    if _, ok := csTerms[c.CourseSectionId]; ok {
      errs = append(errs, errors.New("calendarsessionsection dup entry"))
    }
    csTerms[c.CourseSectionId] = c.Id
  }

  for _, cs := range t.courseSection {
    lingkCourse, ok := courses[cs.CourseId]
    if !ok {
      errs = append(errs, errors.New("missing course id"))
      continue
    }

    termId, ok := csTerms[cs.Id]
    if !ok {
      errs = append(errs, errors.New("coursesection has no calendarsessionsection entry"))
      continue
    }

    if _, ok := p.Courses[cs.CourseId]; !ok {
      p.Courses[cs.CourseId] = &data.Course{
        Title: lingkCourse.Title,
        Description: lingkCourse.Description,
        Campus: lingkCourse.Campus,
      }
    }

    if _, ok := p.Terms[termId]; !ok {
      lingkTerm := terms[termId]
      p.Terms[termId] = &data.Term{
        Semester: lingkTerm.Semester,
        Start: lingkTerm.Start,
        End: lingkTerm.End,
      }
    }

    if _, ok = p.CourseSections[cs.Id]; ok {
      errs = append(errs, errors.New("dup.Coursesection id"))
    }

    p.CourseSections[cs.Id] = &data.CourseSection{
      Course: cs.CourseId,
      Term: termId,
      Section: cs.Section,
      Seats: cs.Seats,
      Status: cs.Status,
      QuarterCredits: cs.QuarterCredits,
      Schedule: make([]*data.Schedule, 0, 1),
      Staff: make([]string, 0, 1),
    }
  }

  for _, s := range t.courseSectionSchedule {
    cs, ok := p.CourseSections[s.CourseSectionId]
    if !ok {
      errs = append(errs, errors.New("schedule slot references nonexistent coursesectoin"))
      continue
    }
    cs.Schedule = append(cs.Schedule, &data.Schedule{
      Days: s.Days,
      Start: s.Start,
      End: s.End,
      Location: s.Location,
    })
  }

  for _, s := range t.sectionInstructor {
    cs, ok := p.CourseSections[s.CourseSectionId]
    if !ok {
      errs = append(errs, errors.New("sectioninstructor references nonexistent coursesection"))
      continue
    }
    st, ok := staff[s.StaffId]
    if !ok {
      errs = append(errs, errors.New("sectioninstructor refs nonexistent staff id"))
      continue
    }

    if _, ok := p.Staff[st.Id]; !ok {
      p.Staff[st.Id] = data.Name{
        First: st.First,
        Last: st.Last,
      }
    }
    cs.Staff = append(cs.Staff, st.Id)
  }

  return p, errs
}
