package fakecsv

import (
	"errors"
	"io"
	"log"
)

type state int

const (
	expectStart state = iota
	inString
	gotQuote
)

type FakeCsvReader struct {
	br io.ByteReader
	st state
}

func New(br io.ByteReader) *FakeCsvReader {
	return &FakeCsvReader{br, expectStart}
}

func (r *FakeCsvReader) ReadRow() ([]string, error) {
	var row []string
	var cell []byte
	for {
		c, err := r.br.ReadByte()
		if err == io.EOF {
		  switch r.st {
		  case gotQuote:
			return append(row, string(cell)), nil
		      case inString:
			return nil, errors.New("failed while reading string")
		      }
		}
		if err != nil {
			return nil, err
		}
		switch r.st {
		case expectStart:
			if c != '"' {
				log.Fatalf("expecting '\"' marking string start; got character '%c' instead\n", c)
			}
			r.st = inString
		case inString:
			switch c {
			case '\\':
				continue
			case '"':
				r.st = gotQuote
				continue
			default:
				cell = append(cell, c)
			}
		case gotQuote:
			switch c {
			case '\n':
				return append(row, string(cell)), nil
			case ',':
				row = append(row, string(cell))
				cell = nil
				r.st = expectStart
			case '"':
				cell = append(cell, '"')
				r.st = gotQuote
			default:
				cell = append(cell, '"')
				r.st = inString
			}

		}
	}
}
