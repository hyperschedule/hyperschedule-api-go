package lingk

import (
	"archive/zip"
	"errors"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsession"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsessionsection"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/course"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesection"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesectionschedule"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/sectioninstructor"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/staff"
	"io"
	"io/ioutil"
	"mime/multipart"
)

func (t *tables) unpackCourse(r io.Reader) error {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	entries, _, _, err := course.ParseAll(contents)
	if err != nil {
		return err
	}
	//allErrs = append(allErrs, errs...)
	t.course = entries
	return nil
}

func (t *tables) unpackCourseSection(r io.Reader) error {
	entries, errs, err := coursesection.ReadAll(r)
	t.warnings = append(t.warnings, errs...)
	t.courseSection = entries
	return err
}

func (t *tables) unpackCourseSectionSchedule(r io.Reader) error {
	entries, errs, err := coursesectionschedule.ReadAll(r)
	t.warnings = append(t.warnings, errs...)
	t.courseSectionSchedule = entries
	return err
}

func (t *tables) unpackCalendarSession(r io.Reader) error {
	entries, errs, err := calendarsession.ReadAll(r)
	t.warnings = append(t.warnings, errs...)
	t.calendarSession = entries
	return err
}

func (t *tables) unpackCalendarSessionSection(r io.Reader) error {
	entries, errs, err := calendarsessionsection.ReadAll(r)
	t.warnings = append(t.warnings, errs...)
	t.calendarSessionSection = entries
	return err
}

func (t *tables) unpackSectionInstructor(r io.Reader) error {
	entries, errs, err := sectioninstructor.ReadAll(r)
	t.warnings = append(t.warnings, errs...)
	t.sectionInstructor = entries
	return err
}

func (t *tables) unpackStaff(r io.Reader) error {
	entries, errs, err := staff.ReadAll(r)
	t.warnings = append(t.warnings, errs...)
	t.staff = entries
	return err
}

func unpackIgnore(_ io.Reader) error { return nil }

func (t *tables) unpackers() map[string]func(r io.Reader) error {
	return map[string]func(r io.Reader) error{
		"course_1.csv":                 t.unpackCourse,
		"coursesection_1.csv":          t.unpackCourseSection,
		"coursesectionschedule_1.csv":  t.unpackCourseSectionSchedule,
		"calendarsession_1.csv":        t.unpackCalendarSession,
		"calendarsessionsection_1.csv": t.unpackCalendarSessionSection,
		"sectioninstructor_1.csv":      t.unpackSectionInstructor,
		"staff_1.csv":                  t.unpackStaff,
		"department_1.csv":             unpackIgnore,
		"departmentcourse_1.csv":       unpackIgnore,
		"facility_1.csv":               unpackIgnore,
		"institution_1.csv":            unpackIgnore,
		"program_1.csv":                unpackIgnore,
		"programcourse_1.csv":          unpackIgnore,
	}
}


func unpack(fh *multipart.FileHeader) (*tables, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := zip.NewReader(f, fh.Size)
	if err != nil {
		return nil, err
	}

	t := &tables{}
	unpackers := t.unpackers()
	for _, mem := range r.File {
    unpacker, ok := unpackers[mem.Name]
    if !ok {
      return nil, errors.New("unrecognized filename")
    }
    r, err := mem.Open()
    if err != nil {
      return nil, errors.New("failed to open zip")
    }
    err = unpacker(r)
    if err != nil {
      return nil, errors.New("failed to unpack/parse csv")
    }
	}

  return t, nil
}
