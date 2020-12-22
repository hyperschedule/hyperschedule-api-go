package course

import "fmt"

type ErrIncorrectHead string

func (s ErrIncorrectHead) Error() string {
	return fmt.Sprintf("Incorrect header: expecting %#v but got %#v", expectHead, s)
}
