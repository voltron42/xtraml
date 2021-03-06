package main

import (
	"../"
	"encoding/xml"
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
	bytes, err = xml.MarshalIndent(q, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

type Node interface {
	Text() string
	Yes() Node
	No() Node
}

type NodeWrapper struct {
	Node Node
}

var choiceParser = choice.ChoiceParser{
	"q": func() interface{} {
		return &Question{}
	},
	"a": func() interface{} {
		return &Answer{}
	},
}

func (n *NodeWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var node *Node
	return choiceParser.ParseList(d, start, &n.Node, node, choice.Set)
}

type Question struct {
	XMLName   xml.Name    `xml:"q"`
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
	XMLName   xml.Name `xml:"a"`
	TextField string   `xml:"text,attr"`
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
