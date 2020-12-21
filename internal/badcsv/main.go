package badcsv

import (
	"fmt"
	"log"
	"regexp"
)

type Course struct {
	dept   string
	num    string
	deptId string
	campus string
	title  string
	desc   string
}

const expectHead = "" +
	`"externalId",` +
	`"classificationOfInstructionalProgramCode",` +
	`"courseNumber",` +
	`"courseTitle",` +
	`"subjectAbbreviation",` +
	`"departmentExternalId",` +
	`"institutionExternalId",` +
	`"facilityExternalId",` +
	`"description"` + "\n"

var reHead = regexp.MustCompile(`^.*\n`)
var reStart = regexp.MustCompile(`(?m)^"([^"]+)","`)
var reRow = regexp.MustCompile(fmt.Sprintf(
	`^"%s","%s","%s","%s","%s","%s","%s","%s","%s"\n?$`,
	`([^"]+)`,    // externalId
	`([^"]*)`,    // classificationOfInstructionalProgramCode
	`([^"]+)`,    // courseNumber (same as externalId)
	`(.*?)`,      // courseTitle
	`([A-Z]*)`,   // subjectAbbreviation
	`([A-Z]*)`,   // departmentExternalId (same as subjectAbbreviation)
	`([A-Z]{2})`, // institutionExternalId
	`([A-Z]{2})`, // facilityExternalId (same as institutionExternalId)
	`((?s).*?)`,  // description
))
var reCode = regexp.MustCompile(`^([A-Z_/-]+) *([0-9A-Z/ -]*?) *([A-Z]{2})?$`)

type ErrIncorrectHead string

func (s ErrIncorrectHead) Error() string {
	return fmt.Sprintf("Incorrect header: expecting %#v but got %#v", expectHead, s)
}

type fail struct {
	full string
	id   string
	data failData
}

type warn struct {
	full string
	id   string
	data warnData
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

type warnData interface {
}
type warnCodeCampusEmpty struct {
	externalId string
}

func Parse(input []byte) ([]*Course, []*warn, []*fail, error) {
	matchHead := string(reHead.Find(input))
	if matchHead != expectHead {
		return nil, nil, nil, ErrIncorrectHead(matchHead)
	}
	input = input[len(matchHead):]

	starts := make([]int, 0, 1024)
	for _, match := range reStart.FindAllIndex(input, -1) {
		starts = append(starts, match[0])
	}

	entries := make([]*Course, 0, len(starts)-1)
	fails := make([]*fail, 0, 8)
	warns := make([]*warn, 0, 8)
	for i, start := range starts[1:] {
		row := string(input[starts[i]:start])
		match := reRow.FindStringSubmatch(row)
		if match == nil {
			log.Printf("failed to parse %#v", row)
			fails = append(fails, &fail{
				full: row,
				data: failFullMatch{},
			})
			continue
		}

		colExternalId := match[1]
		colCourseNumber := match[3]
		colCourseTitle := match[4]
		colSubjectAbbreviation := match[5]
		colDepartmentExternalId := match[6]
		colInstitutionExternalId := match[7]
		colFacilityExternalId := match[8]
		colDescription := match[8]

		if colExternalId != colCourseNumber {
			fails = append(fails, &fail{
				full: row,
				id:   colExternalId,
				data: &failCodeConsistent{
					externalId:   colExternalId,
					courseNumber: colCourseNumber,
				},
			})
			continue
		}

		if colSubjectAbbreviation != colDepartmentExternalId {
			fails = append(fails, &fail{
				full: row,
				id:   colExternalId,
				data: &failSubjectConsistent{
					subjectAbbreviation:  colSubjectAbbreviation,
					departmentExternalId: colDepartmentExternalId,
				},
			})
			continue
		}

		if colInstitutionExternalId != colFacilityExternalId {
			fails = append(fails, &fail{
				full: row,
				id:   colExternalId,
				data: &failCampusConsistent{
					institutionExternalId: colInstitutionExternalId,
					facilityExternalId:    colFacilityExternalId,
				},
			})
			continue
		}

		matchCode := reCode.FindStringSubmatch(colExternalId)
		if matchCode == nil {
			fails = append(fails, &fail{
				full: row,
				id:   colExternalId,
				data: &failCodeMatch{
					input: match[1],
				},
			})
			continue
		}

		codeDept := matchCode[1]
		codeNum := matchCode[2]
		codeCampus := matchCode[3]

		if len(codeCampus) == 0 {
			warns = append(warns, &warn{
				full: row,
				id:   colExternalId,
				data: &warnCodeCampusEmpty{
					externalId: match[1],
				},
			})
		}

		entries = append(entries, &Course{
			dept:   codeDept,
			num:    codeNum,
			deptId: colDepartmentExternalId,
			campus: colInstitutionExternalId,
			title:  colCourseTitle,
			desc:   colDescription,
		})
	}
	return entries, warns, fails, nil
}
