package data

import (
	"fmt"
	"sort"
)

type CourseKey struct {
	Department string
	Code       string
	Campus     string
}

func (k CourseKey) String() string {
	return fmt.Sprintf("%s %s %s", k.Department, k.Code, k.Campus)
}

type SectionKey struct {
	Course  CourseKey
	Term    string
	Section int
}

func (k SectionKey) String() string {
	return fmt.Sprintf("%s-%02d %s", k.Course, k.Section, k.Term)
}

type Data struct {
	CourseSections map[SectionKey]*CourseSection
	Courses        map[CourseKey]*Course
	Terms          map[string]*Term
	Staff          map[string]Name
}

type CourseSection struct {
	Course         CourseKey
	Term           string
	Section        int
	Seats          Seats
	Status         Status
	QuarterCredits int
	Schedule       []*Schedule
	Staff          []string
}

type Course struct {
	Name        string
	Description string
}

type Term struct {
	Semester string
	Start    Date
	End      Date
}

type Name struct {
	First string
	Last  string
	Alt   *string
}

type Schedule struct {
	Days     Days
	Start    Time
	End      Time
	Location string
}

func ScheduleLess(s1, s2 *Schedule) bool {
	return s1.Days < s2.Days || s1.Days == s2.Days &&
		(TimeLess(s1.Start, s2.Start) || s1.Start == s2.Start &&
			(TimeLess(s1.End, s2.End) || s1.End == s2.End &&
				(s1.Location < s2.Location)))
}

func (s *CourseSection) ScheduleLess(i, j int) bool {
	return ScheduleLess(s.Schedule[i], s.Schedule[j])
}

func SectionsEqual(s1 *CourseSection, s2 *CourseSection) bool {
	if len(s1.Schedule) != len(s2.Schedule) || len(s1.Staff) != len(s2.Staff) {
		return false
	}

	// TODO use maps/sets to compare equality

	sort.StringSlice(s1.Staff).Sort()
	sort.StringSlice(s2.Staff).Sort()

	for i := range s1.Staff {
		if s1.Staff[i] != s2.Staff[i] {
			return false
		}
	}

	sort.Slice(s1.Schedule, s1.ScheduleLess)
	sort.Slice(s2.Schedule, s2.ScheduleLess)

	for i := range s1.Schedule {
		if s1.Schedule[i] != s2.Schedule[i] {
			return false
		}
	}

	return s1.Course == s2.Course &&
		s1.Term == s2.Term &&
		s1.Section == s2.Section &&
		s1.Seats == s2.Seats &&
		s1.Status == s2.Status &&
		s1.QuarterCredits == s2.QuarterCredits
}
