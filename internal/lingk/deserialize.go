package lingk

import (
	"bytes"
	"embed"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	_ "log"
	"mime/multipart"
)

//go:embed sample
var sample embed.FS

func FromAttachment(fh *multipart.FileHeader) (*data.Data, error) {
	t, err := Unpack(fh)
	if err != nil {
		return nil, err
	}

	d, errs := t.prune()
	_ = errs
	//for _, err := range errs {
	//	log.Printf("warning: %v", err)
	//}

	return d, nil
}

func FromZipBuffer(buf *bytes.Buffer) (*data.Data, error) {
	t, err := UnpackZip(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return nil, err
	}

	d, errs := t.prune()
	_ = errs
	//for _, err := range errs {
	//	log.Printf("warning: %v", err)
	//}

	return d, nil
}

func Sample() (*data.Data, error) {
	t, err := unpackFs(sample, "sample/fa2021")
	if err != nil {
		return nil, err
	}

	d, _ := t.prune()
	return d, nil
}
