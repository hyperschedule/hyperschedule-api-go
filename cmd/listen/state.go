package main

// TODO thread safety

import (
  "github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

type State struct{
  data *data.Data
}

var state State

func (s *State) SetData(d *data.Data) {
  s.data = d
}

func (s *State) GetData() *data.Data {
  return s.data
}
