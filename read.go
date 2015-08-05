package xtraml

import (
	"errors"
	"io"
	"reflect"
)

func ReadXML(reader io.Reader, ptr interface{}) error {
	start, err := first(reader)
	if err != nil {
		return err
	}
	value := reflect.ValueOf(ptr).Elem()
	return read(*start, value)
}

func read(start StartToken, value reflect.Value) error {
	obj := value.Interface()
	readable, ok := obj.(Readable)
	if !ok {
		if value.Kind() != reflect.Struct {
			return errors.New("Cannot read to an object which is not a struct")
		}
		v := valueReadable{value}
		readable = &v
	}
	return readable.ReadXML(start)
}

type Readable interface {
	ReadXML(start StartToken) error
}

type ReadableAttr interface {
	ReadAttr(name Name, value string) error
}

type valueReadable struct {
	value reflect.Value
}

func (v *valueReadable) ReadXML(start StartToken) error {
	// TODO:
	return nil
}
