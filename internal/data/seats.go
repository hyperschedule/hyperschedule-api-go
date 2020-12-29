package data

import (
  "strconv"
  "errors"
)

type Seats struct{
  Enrolled int
  Capacity int
}

func ParseSeats(enrolled string, capacity string) (Seats, error) {
  nCapacity, err := strconv.Atoi(capacity)
  if err != nil {
    return Seats{}, errors.New("invalid capacity")
  }

  nEnrolled, err := strconv.Atoi(enrolled)
  if err != nil {
    return Seats{}, errors.New("invalid enrollment")
  }

  return Seats{Enrolled: nEnrolled, Capacity: nCapacity}, nil
}
