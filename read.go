package xtraml

import (
	"io"
)

func ReadXML(reader io.Reader, ptr interface{}) error {
	start, err := first(reader)
	
	return nil
}

type Readable interface {
	ReadXML(start StartToken) error
}
