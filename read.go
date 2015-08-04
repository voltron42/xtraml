package xtraml

import (
	"io"
)

func ReadXML(reader io.Reader, ptr interface{}) error {
	// TODO:
	return nil
}

type Readable interface {
	ReadXML(start StartToken) error
}
