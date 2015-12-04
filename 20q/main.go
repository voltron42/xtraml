package main

import (
	"../choice"
	"encoding/xml"
	"errors"
	"fmt"
)

func main() {
	bytes := []byte("<q text=\"q1\"><yes><a text=\"a1\"/></yes><no><q text=\"q2\"><yes><a text=\"a2\"/></yes><no><a text=\"a3\"/></no></q></no></q>")
	q := Question{}
	err := xml.Unmarshal(bytes, &q)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", q)
}

type Node interface {
	Text() string
	Yes() Node
	No() Node
}

type NodeWrapper struct {
	Node
}

var choiceParser = choice.ChoiceParser{
	"q": func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		q := Question{}
		err := d.DecodeElement(&q, &start)
		return q, err
	},
	"a": func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		a := Answer{}
		err := d.DecodeElement(&a, &start)
		return a, err
	},
}

func (n *NodeWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return choiceParser.ParseList(d, start, func(item interface{}) error {
		node, ok := item.(Node)
		if !ok {
			return errors.New("Item is not a valid Node.")
		}
		n.Node = node
		return nil
	})
}

type Question struct {
	TextField string      `xml:"text,attr"`
	YesField  NodeWrapper `xml:"yes"`
	NoField   NodeWrapper `xml:"no"`
}

func (q Question) Text() string {
	return q.TextField
}

func (q Question) Yes() Node {
	return q.YesField.Node
}

func (q Question) No() Node {
	return q.NoField.Node
}

type Answer struct {
	TextField string `xml:"text,attr"`
}

func (a Answer) Text() string {
	return a.TextField
}

func (a Answer) Yes() Node {
	return nil
}

func (a Answer) No() Node {
	return nil
}
