package main

import (
	"encoding/xml"
	"fmt"
)

func main() {
	bytes := []byte("<q text=\"q1\"><yes><a text=\"a1\"/></yes><no><q test=\"q1\"><yes><a text=\"a2\"/></yes><no><a text=\"a3\"/></no></q></no></q>")
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

func (n NodeWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fmt.Printf("name: %v:%v\n", start.Name.Space, start.Name.Local)
	token, err := d.Token()
	if err != nil {
		return err
	}
	fmt.Printf("token: %v", token)
	node, ok := token.(xml.StartElement)
	if !ok {
		return errors.New("Wrapped Node contains no elements")
	}
	return nil
}
func (n NodeWrapper) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(n.Node, start)
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
