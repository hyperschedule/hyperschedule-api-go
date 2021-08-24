package main

// TODO thread safety

import (
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

type OldState struct {
	data *data.Data
}

func (s *OldState) SetData(d *data.Data) {
	s.data = d
}

func (s *OldState) GetData() *data.Data {
	return s.data
}
