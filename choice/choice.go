package choice

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
)

type appender int

const (
	Set appender = iota
	Append
)

var appenders = map[appender]func(container reflect.Value, value reflect.Value){
	Set: func(container reflect.Value, value reflect.Value) {
		container.Set(value)
	},
	Append: func(container reflect.Value, value reflect.Value) {
		container.Set(reflect.Append(container, value))
	},
}

type ChoiceParser map[string]func() interface{}

func (c ChoiceParser) Parse(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
	choice, ok := c[start.Name.Local]
	if !ok {
		return nil, errors.New(start.Name.Local + " is not a listed option.")
	}
	ptr := choice()
	err := d.DecodeElement(ptr, &start)
	obj := reflect.ValueOf(ptr).Interface()
	return obj, err
}

func (c ChoiceParser) ParseList(d *xml.Decoder, start xml.StartElement, containerPtr interface{}, typeofPtr interface{}, appenderType appender) error {
	typeof := reflect.TypeOf(typeofPtr).Elem()
	container := reflect.ValueOf(containerPtr).Elem()
	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}
		next, ok := token.(xml.StartElement)
		if ok {
			item, err := c.Parse(d, next)
			if err != nil {
				return err
			}
			val := reflect.ValueOf(item)
			if !val.Type().Implements(typeof) {
				return fmt.Errorf("Item is not a valid %v.", typeof.Name())
			}
			appendFn := appenders[appenderType]
			appendFn(container, val)
		}
		token, err = d.Token()
	}
	return nil
}

func WrapList(e *xml.Encoder, listName xml.Name, list interface{}) error {
	listVal := reflect.ValueOf(list)
	token := xml.StartElement{
		Name: listName,
	}
	err := e.EncodeToken(token)
	if err != nil {
		return err
	}
	count := listVal.Len()
	for index := 0; index < count; index++ {
		err = e.Encode(listVal.Index(index).Interface())
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(token.End())

}
