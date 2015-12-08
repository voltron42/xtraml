package main

import (
	"../"
	"encoding/xml"
	"fmt"
)

func main() {
	sample := "<a><list><b value=\"this is b\"/><c>this is c</c><d><b value=\"this is b\"/><c>this is c</c></d></list></a>"
	a := A{}
	err := xml.Unmarshal([]byte(sample), &a)
	if err != nil {
		panic(err)
	}
	for _, item := range a.List {
		fmt.Printf("%v\n", item.ToString())
	}
	fmt.Printf("%v\n", a)
	bytes, err := xml.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

type A struct {
	XMLName xml.Name `xml:"a"`
	List    List     `xml:"list"`
}

type List []Item

var itemChoices = choice.ChoiceParser{
	"b": func() interface{} {
		return &B{}
	},
	"c": func() interface{} {
		return &C{}
	},
	"d": func() interface{} {
		return &D{}
	},
}

func (p *List) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var item *Item
	return itemChoices.ParseList(d, start, p, item, choice.Append)
}

func (p List) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return choice.WrapList(e, xml.Name{Local: "list"}, p)
}

type Item interface {
	ToString() string
}

type B struct {
	XMLName xml.Name `xml:"b"`
	Value   string   `xml:"value,attr"`
}

func (b B) ToString() string {
	return b.Value
}

type C struct {
	XMLName xml.Name `xml:"c"`
	Data    string   `xml:",chardata"`
}

func (c C) ToString() string {
	return c.Data
}

type D struct {
	XMLName xml.Name `xml:"d"`
	B       B        `xml:"b"`
	C       C        `xml:"c"`
}

func (d D) ToString() string {
	return fmt.Sprintf("b=(%v)\nc=(%v)\n", d.B.ToString(), d.C.ToString())
}
