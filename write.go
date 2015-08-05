package xtraml

import (
	"errors"
	"io"
	"reflect"
)

func WriteXML(obj interface{}, writer io.Writer) error {
	writable, ok := obj.(Writeable)
	if !ok {
		w, err := newValueWriteable(obj)
		if err != nil {
			return err
		}
		writable = w
	}
	return writable.WriteXML(writer)
}

func newValueWriteable(obj interface{}) (Writeable, error) {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Struct {
		return nil, errors.New("Cannot write from an object which is not a Struct")
	}
	return &valueWriteable{value}, nil
}

type Writeable interface {
	WriteXML(writer io.Writer) error
}

type valueWriteable struct {
	value reflect.Value
}

func (v valueWriteable) WriteXML(writer io.Writer) error {
	// TODO:
	return nil
}
