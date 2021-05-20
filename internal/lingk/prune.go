package lingk

import (
	"errors"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsession"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/course"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/staff"
)

func (t *tables) prune() (*data.Data, []error) {
	p := &data.Data{
		Terms:          make(map[string]*data.Term),
		Courses:        make(map[data.CourseKey]*data.Course),
		CourseSections: make(map[data.SectionKey]*data.CourseSection),
		Staff:          make(map[string]data.Name),
	}
	sectionIds := map[string]struct{}{}
	errs := make([]error, 0)

	// We don't add directly to `p.Courses` here because the raw `course_1.csv`
	// table contains all sorts of extraneous (bad) entries that are never
	// actually referenced by the `coursesection_1.csv` entries--which are the
	// items we actually care about, and which are far more
	// well-formed/consistent than the `course_1.csv` entries.

	// So we build an intermediate dictionary `courses` below, but only
	// ("lazily") add them to our final dictionary `p.Courses` when they _do_ get
	// referenced by a `coursesection` entry.

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

	csSchedule := make(map[string][]*data.Schedule)
	for _, c := range t.courseSectionSchedule {
		csSchedule[c.CourseSectionId] = append(csSchedule[c.CourseSectionId], &data.Schedule{
			Days:     c.Days,
			Start:    c.Start,
			End:      c.End,
			Location: c.Location,
		})
	}

	csStaff := make(map[string][]string)
	for _, s := range t.sectionInstructor {
		st, ok := staff[s.StaffId]
		if !ok {
			errs = append(errs, errors.New("sectioninstructor refs nonexistent staff id"))
			continue
		}

		if _, ok := p.Staff[s.StaffId]; !ok {
			p.Staff[s.StaffId] = data.Name{First: st.First, Last: st.Last}
		}

		csStaff[s.CourseSectionId] = append(csStaff[s.CourseSectionId], s.StaffId)
	}

	for _, cs := range t.courseSection {
		lingkCourse, ok := courses[cs.CourseId]
		if !ok {
			errs = append(errs, errors.New("coursesection has no matching course"))
			continue
		}

		termId, ok := csTerms[cs.Id]
		if !ok {
			errs = append(errs, errors.New("coursesection has no calendarsessionsection entry"))
			continue
		}

		courseKey := data.CourseKey{
			Department: lingkCourse.Department,
			Code:       lingkCourse.Number,
			Campus:     lingkCourse.Campus,
		}
		if _, ok := p.Courses[courseKey]; !ok {
			p.Courses[courseKey] = &data.Course{
				Name:        lingkCourse.Title,
				Description: lingkCourse.Description,
			}
		}

		if _, ok := p.Terms[termId]; !ok {
			lingkTerm := terms[termId]
			p.Terms[termId] = &data.Term{
				Semester: lingkTerm.Semester,
				Start:    lingkTerm.Start,
				End:      lingkTerm.End,
			}
		}

		sectionKey := data.SectionKey{
			Course:  courseKey,
			Term:    termId,
			Section: cs.Section,
		}
		if _, ok := p.CourseSections[sectionKey]; ok {
			errs = append(errs, errors.New("dup coursesection key"))
		}

		if _, ok := sectionIds[cs.Id]; ok {
			errs = append(errs, errors.New("dup coursesection id"))
		}

		p.CourseSections[sectionKey] = &data.CourseSection{
			Course:         courseKey,
			Term:           termId,
			Section:        cs.Section,
			Seats:          cs.Seats,
			Status:         cs.Status,
			QuarterCredits: cs.QuarterCredits,
			Schedule:       csSchedule[cs.Id],
			Staff:          csStaff[cs.Id],
		}
	}

	return p, errs
}
