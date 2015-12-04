package main

import (
	"../choice"
	"encoding/xml"
	"errors"
	"fmt"
)

func main() {
	sample := "<a><b value=\"this is b\"/><c>this is c</c><d><b value=\"this is b\"/><c>this is c</c></d></a>"
	list := List{}
	err := xml.Unmarshal([]byte(sample), &list)
	if err != nil {
		panic(err)
	}
	for _, item := range list {
		fmt.Printf("%v\n", item.ToString())
	}
	fmt.Printf("%v\n", list)
}

type List []Item

var itemChoices = choice.ChoiceParser{
	"b": func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		a := B{}
		err := d.DecodeElement(&a, &start)
		return a, err
	},
	"c": func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		a := C{}
		err := d.DecodeElement(&a, &start)
		return a, err
	},
	"d": func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		a := D{}
		err := d.DecodeElement(&a, &start)
		return a, err
	},
}

func (p *List) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return itemChoices.ParseList(d, start, func(item interface{}) error {
		pathItem, ok := item.(Item)
		if !ok {
			return errors.New("Item is not a valid PathItem.")
		}
		*p = append(*p, pathItem)
		return nil
	})
}

type Item interface {
	ToString() string
}

type B struct {
	Value string `xml:"value,attr"`
}

func (b B) ToString() string {
	return b.Value
}

type C struct {
	Data string `xml:",chardata"`
}

func (c C) ToString() string {
	return c.Data
}

type D struct {
	B B `xml:"b"`
	C C `xml:"c"`
}

func (d D) ToString() string {
	return fmt.Sprintf("b=(%v)\nc=(%v)\n", d.B.ToString(), d.C.ToString())
}
