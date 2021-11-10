package update

import (
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

type rows struct {
	term            [][]interface{}
	staff           [][]interface{}
	course          [][]interface{}
	section         [][]interface{}
	sectionStaff    [][]interface{}
	sectionSchedule [][]interface{}
}

func (rs rows) seq() [6][][]interface{} {
	return [...][][]interface{}{
		rs.term,
		rs.staff,
		rs.course,
		rs.section,
		rs.sectionStaff,
		rs.sectionSchedule,
	}
}

func rowsFrom(d *data.Data) *rows {
	rs := &rows{
		staff:           make([][]interface{}, 0, len(d.Staff)),
		term:            make([][]interface{}, 0, len(d.Terms)),
		course:          make([][]interface{}, 0, len(d.Courses)),
		section:         make([][]interface{}, 0, len(d.CourseSections)),
		sectionStaff:    make([][]interface{}, 0, len(d.CourseSections)),
		sectionSchedule: make([][]interface{}, 0, len(d.CourseSections)),
	}

	for lingkId, name := range d.Staff {
		rs.staff = append(rs.staff, []interface{}{lingkId, name.First, name.Last, name.Alt})
	}

	for code, term := range d.Terms {
		rs.term = append(rs.term, []interface{}{
			code, term.Semester, term.Start.ToTime(), term.End.ToTime(),
		})
	}

	for key, course := range d.Courses {
		rs.course = append(rs.course, []interface{}{
			key.Department,
			key.Code,
			key.Campus,
			course.Name,
			course.Description,
		})
	}

	for key, section := range d.CourseSections {
		keyCols := []interface{}{
			key.Course.Department,
			key.Course.Code,
			key.Course.Campus,
			key.Term,
			key.Section,
		}

		rs.section = append(rs.section, append(keyCols,
			section.QuarterCredits,
			section.Status.String(),
			section.Seats.Enrolled,
			section.Seats.Capacity,
		))

		for _, lingkId := range section.Staff {
			rs.sectionStaff = append(rs.sectionStaff, append(keyCols, lingkId))
		}

		for _, schedule := range section.Schedule {
			rs.sectionSchedule = append(rs.sectionSchedule, append(keyCols,
				int(schedule.Days),
				schedule.Start.Std(),
				schedule.End.Std(),
				schedule.Location,
			))
		}
	}

	return rs
}
