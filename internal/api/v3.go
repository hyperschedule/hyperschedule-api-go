package api

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"sort"
	"strings"
	"time"
)

type V3 struct {
	Data  *V3CourseData `json:"data"`
	Until int64         `json:"until"`
	Error *string       `json:"error"`
	Full  bool          `json:"full"`
}

type V3CourseData struct {
	Terms   map[string]*V3Term   `json:"terms"`
	Courses map[string]*V3Course `json:"courses"`
}

type V3Term struct {
	Code    string        `json:"termCode"`
	SortKey []interface{} `json:"termSortKey"`
	Name    string        `json:"termName"`
}

type V3Course struct {
	Code               string        `json:"courseCode"`
	Name               string        `json:"courseName"`
	SortKey            []interface{} `json:"courseSortKey"`
	MutualExclusionKey []interface{} `json:"courseMutualExclusionKey"`
	Description        string        `json:"courseDescription"`
	Instructors        []string      `json:"courseInstructors"`
	Term               string        `json:"courseTerm"`
	Schedule           []*V3Schedule `json:"courseSchedule"`
	Credits            float32       `json:"courseCredits"`
	SeatsTotal         int           `json:"courseSeatsTotal"`
	SeatsFilled        int           `json:"courseSeatsFilled"`
	WaitlistLength     *int          `json:"courseWaitlistLength"`
	EnrollmentStatus   string        `json:"courseEnrollmentStatus"`
	PermCount          int           `json:"permCount"`
}

type V3Schedule struct {
	Days      string `json:"scheduleDays"`
	StartTime string `json:"scheduleStartTime"`
	EndTime   string `json:"scheduleEndTime"`
	StartDate string `json:"scheduleStartDate"`
	EndDate   string `json:"scheduleEndDate"`
	TermCount int    `json:"scheduleTermCount"`
	Terms     []int  `json:"scheduleTerms"`
	Location  string `json:"scheduleLocation"`
}

func MakeV3(d *data.Data) *V3 {
	terms := make(map[string]*V3Term)
	for id, term := range d.Terms {
		terms[id] = TermToV3Term(id, term)
	}

	return &V3{
		Data: &V3CourseData{
			Terms:   terms,
			Courses: MakeV3Courses(d),
		},
		Until: time.Now().Unix(),
		Error: nil,
		Full:  true,
	}
}

func TermToV3Term(id string, t *data.Term) *V3Term {
	return &V3Term{
		Code:    id,
		Name:    t.Semester,
		SortKey: []interface{}{id, t.Semester},
	}
}

func MakeV3Courses(d *data.Data) map[string]*V3Course {
	courses := make(map[string]*V3Course)
	for id, cs := range d.CourseSections {
		course := d.Courses[cs.Course]
		term := d.Terms[cs.Term]

		// more dirty hacks
		if !strings.HasPrefix(cs.Term, "FA2021") {
			continue
		}

		code := fmt.Sprintf("%s %s %s-%02d",
			cs.Course.Department,
			cs.Course.Code,
			cs.Course.Campus,
			cs.Section,
		)

		instructors := make([]string, len(cs.Staff))
		for i, staff := range cs.Staff {
			name := d.Staff[staff]
			instructors[i] = fmt.Sprintf("%s %s", name.First, name.Last)
		}
		sort.StringSlice(instructors).Sort()

		// these are dirty hacks, fix
		termCount := 1
		if cs.Term != term.Semester {
			termCount = 2
		}

		termIndex := 0
		if termCount == 2 && strings.HasSuffix(cs.Term, "2") {
			termIndex = 1
		}

		schedule := make([]*V3Schedule, 0)
		for _, sched := range cs.Schedule {

			schedule = append(schedule, &V3Schedule{
				Days:      sched.Days.String(),
				StartTime: sched.Start.String(),
				EndTime:   sched.End.String(),
				StartDate: term.Start.String(),
				EndDate:   term.End.String(),
				TermCount: termCount,
				Terms:     []int{termIndex},
				Location:  sched.Location,
			})
		}
		sort.Slice(schedule, func(i, j int) bool {
			return v3ScheduleLess(schedule[i], schedule[j])
		})

		courses[code] = &V3Course{
			Code: code,
			Name: course.Name,
			SortKey: []interface{}{
				fmt.Sprintf("%s-%02d", id.Course.String(), id.Section),
			},
			MutualExclusionKey: []interface{}{cs.Course.String()},
			Description:        course.Description,
			Instructors:        instructors,
			Term:               cs.Term,
			Schedule:           schedule,
			Credits:            float32(cs.QuarterCredits) / 4,
			SeatsTotal:         cs.Seats.Capacity,
			SeatsFilled:        cs.Seats.Enrolled,
			WaitlistLength:     nil,
			EnrollmentStatus:   cs.Status.String(),
		}
	}

	return courses
}
