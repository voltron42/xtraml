package xtraml

import (
	"io"
	"reflect"
)

func WriteXML(obj interface{}, writer io.Writer) error {
	writable, ok := obj.(Writeable)
	if !ok {
		writable = newValueWriteable(obj)
	}
	return writable.WriteXML(writer)
}

func newValueWriteable(obj interface{}) Writeable {
	return &ValueWriteable{reflect.ValueOf(obj)}
}

type Writeable interface {
	WriteXML(writer io.Writer) error
}

type ValueWriteable struct {
	value reflect.Value
}

func (v ValueWriteable) WriteXML(writer io.Writer) error {
	// TODO:
	return nil
}
