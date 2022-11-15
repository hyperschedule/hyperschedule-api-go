package update

import (
	_ "embed"
)

var (
	sectionKeyCols = []columnSpec{
		{"course_department", "text NOT NULL"},
		{"course_code", "text NOT NULL"},
		{"course_campus", "text NOT NULL"},
		{"term_code", "text NOT NULL"},
		{"section", "integer NOT NULL"},
	}

	sectionKeyColNames = colNames(sectionKeyCols)

	tmpTables = tables{

		term: table{
			name: "term_tmp",
			columns: []columnSpec{
				{"code", "text PRIMARY KEY"},
				{"semester", "text NOT NULL"},
				{"date_start", "date NOT NULL"},
				{"date_end", "date NOT NULL"},
			},
		},

		staff: table{
			name: "staff_tmp",
			columns: []columnSpec{
				{"lingk_id", "text PRIMARY KEY"},
				{"name_first", "text NOT NULL"},
				{"name_last", "text NOT NULL"},
				{"alt", "text"},
			},
		},

		course: table{
			name: "course_tmp",
			columns: []columnSpec{
				{"department", "text NOT NULL"},
				{"code", "text NOT NULL"},
				{"campus", "text NOT NULL"},
				{"name", "text NOT NULL"},
				{"description", "text NOT NULL"},
			},
			constraints: []constraintSpec{pkey("department", "code", "campus")},
		},

		section: table{
			name: "section_tmp",
			columns: append(sectionKeyCols, []columnSpec{
				{"quarter_credits", "integer NOT NULL"},
				{"status", "status  NOT NULL"},
				{"seats_enrolled", "integer NOT NULL"},
				{"seats_capacity", "integer NOT NULL"},
				{"perms", "integer NOT NULL"},
			}...),
			constraints: []constraintSpec{pkey(sectionKeyColNames...)},
		},

		sectionStaff: table{
			name: "section_staff_tmp",
			columns: append(sectionKeyCols,
				columnSpec{"staff_lingk_id", "text NOT NULL"},
			),
			constraints: []constraintSpec{
				pkey(append(sectionKeyColNames, "staff_lingk_id")...),
			},
		},

		sectionSchedule: table{
			name: "section_schedule_tmp",
			columns: append(sectionKeyCols, []columnSpec{
				{"days", "smallint NOT NULL"},
				{"time_start", "time NOT NULL"},
				{"time_end", "time NOT NULL"},
				{"location", "text NOT NULL"},
			}...),
			constraints: []constraintSpec{pkey(append(sectionKeyColNames,
				"days",
				"time_start",
				"time_end",
				"location",
			)...)},
		},
	}

	//go:embed sql/term-upsert.sql
	sqlTermUpsert string

	//go:embed sql/term-delete.sql
	sqlTermDelete string

	//go:embed sql/staff-upsert.sql
	sqlStaffUpsert string

	//go:embed sql/staff-delete.sql
	sqlStaffDelete string

	//go:embed sql/course-upsert.sql
	sqlCourseUpsert string

	//go:embed sql/course-delete.sql
	sqlCourseDelete string

	//go:embed sql/section-upsert.sql
	sqlSectionUpsert string

	//go:embed sql/section-delete.sql
	sqlSectionDelete string

	//go:embed sql/section_staff-upsert.sql
	sqlSectionStaffUpsert string

	//go:embed sql/section_staff-delete.sql
	sqlSectionStaffDelete string

	//go:embed sql/section_schedule-upsert.sql
	sqlSectionScheduleUpsert string

	//go:embed sql/section_schedule-delete.sql
	sqlSectionScheduleDelete string
)

func init() {
	tmpTables.term.sqlApply = sqlApply{
		upsert: sqlTermUpsert,
		delete: sqlTermDelete,
	}
	tmpTables.staff.sqlApply = sqlApply{
		upsert: sqlStaffUpsert,
		delete: sqlStaffDelete,
	}
	tmpTables.course.sqlApply = sqlApply{
		upsert: sqlCourseUpsert,
		delete: sqlCourseDelete,
	}
	tmpTables.section.sqlApply = sqlApply{
		upsert: sqlSectionUpsert,
		delete: sqlSectionDelete,
	}
	tmpTables.sectionStaff.sqlApply = sqlApply{
		upsert: sqlSectionStaffUpsert,
		delete: sqlSectionStaffDelete,
	}
	tmpTables.sectionSchedule.sqlApply = sqlApply{
		upsert: sqlSectionScheduleUpsert,
		delete: sqlSectionScheduleDelete,
	}
}
