package lingk

import (
  "github.com/MuddCreates/hyperschedule-api-go/internal/data"
  "mime/multipart"
)

func FromAttachment(fh *multipart.FileHeader) (*data.Data, error) {

  t, err := Unpack(fh)
  if err != nil {
    return nil, err
  }

  d, _ := t.prune()

  return d, nil
}
