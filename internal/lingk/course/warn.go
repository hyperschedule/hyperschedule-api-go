package course

type warn struct {
	full string
	id   string
	data warnData
}

type warnData interface {
}
type warnCodeCampusEmpty struct {
	externalId string
}
