package update

type tables struct {
	term            table
	staff           table
	course          table
	section         table
	sectionStaff    table
	sectionSchedule table
}

type table struct {
	name        string
	columns     []columnSpec
	constraints []constraintSpec
	sqlApply    sqlApply
}

type sqlApply struct {
	upsert string
	delete string
}

type namedSql struct {
	name string
	sql  string
}

type columnSpec struct {
	name  string
	props string
}

type constraintSpec string

func colNames(cols []columnSpec) []string {
	names := make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.name
	}
	return names
}

func (tbls tables) seq() [6]table {
	return [...]table{
		tbls.term,
		tbls.staff,
		tbls.course,
		tbls.section,
		tbls.sectionStaff,
		tbls.sectionSchedule,
	}
}
