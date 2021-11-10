package altstaff

import (
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
	"io"
	"regexp"
	"strings"
)

var reAltName = regexp.MustCompile(`^([A-Za-z_. ]+)(?:\\, ([A-Za-z_. ]+))?`)

var expectHead = []string{
	"externalId",
	"firstName",
	"lastName",
	"altName",
}

func parse(record []string) (*Entry, error) {
	match := reAltName.FindStringSubmatch(record[3])
	if match == nil {
		return nil, fmt.Errorf("failed to parse alt name: %#v", record[3])
	}
	alt := strings.TrimSpace(fmt.Sprintf("%s %s", match[2], match[1]))

	return &Entry{
		Id:    record[0],
		First: record[1],
		Last:  record[2],
		Alt:   alt,
	}, nil
}

func ReadAll(r io.Reader) ([]*Entry, []error, error) {
	entries := make([]*Entry, 0, 128)
	errs, err := csvutil.Collect(r, expectHead, func(record []string) error {
		entry, err := parse(record)
		if err != nil {
			return err
		}
		entries = append(entries, entry)
		return nil
	})
	return entries, errs, err
}
