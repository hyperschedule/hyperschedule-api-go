package lingk

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/altstaff"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsession"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/calendarsessionsection"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/course"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesection"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/coursesectionschedule"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/sectioninstructor"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk/staff"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"path"
)

func (t *tables) unpackCourse(r io.Reader) error {
	contents, err := ioutil.ReadAll(transform.NewReader(r, charmap.Windows1252.NewDecoder()))
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

func (t *tables) unpackAltStaff(r io.Reader) error {
	entries, errs, err := altstaff.ReadAll(r)
	t.warnings = append(t.warnings, errs...)
	t.altStaff = entries
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
		"altstaff_1.csv":               t.unpackAltStaff,
		"department_1.csv":             unpackIgnore,
		"departmentcourse_1.csv":       unpackIgnore,
		"facility_1.csv":               unpackIgnore,
		"institution_1.csv":            unpackIgnore,
		"program_1.csv":                unpackIgnore,
		"programcourse_1.csv":          unpackIgnore,
	}
}

func Unpack(fh *multipart.FileHeader) (*tables, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return UnpackZip(f, fh.Size)
}

func UnpackZip(f io.ReaderAt, size int64) (*tables, error) {
	r, err := zip.NewReader(f, size)
	if err != nil {
		return nil, err
	}

	t := &tables{}
	unpackers := t.unpackers()
	for _, mem := range r.File {
		unpacker, ok := unpackers[mem.Name]
		if !ok {
			return nil, errors.New(fmt.Sprintf("unrecognized filename: %#v", mem.Name))
		}
		r, err := mem.Open()
		defer r.Close()
		if err != nil {
			return nil, err
		}
		err = unpacker(r)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func unpackFs(f fs.FS, dir string) (*tables, error) {
	entries, err := fs.ReadDir(f, dir)
	if err != nil {
		return nil, err
	}

	t := &tables{}
	unpackers := t.unpackers()
	for _, mem := range entries {
		unpacker, ok := unpackers[mem.Name()]
		if !ok {
			return nil, errors.New(fmt.Sprintf("unrecognized filename %s", mem.Name()))
		}
		file, err := f.Open(path.Join(dir, mem.Name()))
		if err != nil {
			return nil, err
		}
		unpacker(file)
	}
	return t, nil
}
