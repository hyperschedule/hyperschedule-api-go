package permcount

import (
	"errors"
	"github.com/MuddCreates/hyperschedule-api-go/internal/csvutil"
	"io"
	"strconv"
)

var expectHead = []string{
	"permCountExternalId",
	"PermCount",
}

func parse(record []string) (*Entry, error) {
	colId := record[0]
	colCount := record[1]

	count, err := strconv.Atoi(colCount)
	if err != nil {
		return nil, errors.New("invalid permcount")
	}

	return &Entry{
		Id:    colId,
		Count: count,
	}, nil
}

func ReadAll(r io.Reader) ([]*Entry, []error, error) {
	countEntries := make([]*Entry, 0, 1024)
	errs, err := csvutil.Collect(r, expectHead, func(record []string) error {
		entry, err := parse(record)
		if err != nil {
			return err
		}
		countEntries = append(countEntries, entry)
		return nil
	})
	return countEntries, errs, err
}
