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
	fmt.Printf("start: %v\n", start.Name)
	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}
		next, ok := token.(xml.StartElement)
		if ok {
			fmt.Printf("next: %v\n", next.Name)
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
