package course

type fail struct {
	full string
	id   string
	data failData
}

type failData interface {
}

type failFullMatch struct{}
type failCodeConsistent struct {
	externalId   string
	courseNumber string
}
type failCodeMatch struct {
	input string
}
type failCampusConsistent struct {
	institutionExternalId string
	facilityExternalId    string
}
type failSubjectConsistent struct {
	subjectAbbreviation  string
	departmentExternalId string
}
