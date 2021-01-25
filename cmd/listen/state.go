package main

import (
  "github.com/MuddCreates/hyperschedule-api-go/internal/data"
)

type State struct{
  data *data.Data
}

func (s *State) Listen() {
}

func (s *State) GetData() {
}
