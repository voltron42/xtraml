package xtraml

import (
	"io"
)

type Node interface {
	WriteXML(writer io.Writer) error
}

type NodeObj struct {
	Name     Name
	Attrs    AttrSet
	Children []Node
}

func (n NodeObj) WriteXML(writer io.Writer) error {
	writer.Write([]byte("<"))
	writer.Write([]byte(n.Name.String()))
	if len(n.Children) < 1 {
		writer.Write([]byte("/>"))
	} else {
		writer.Write([]byte(">"))
    for _, child := range n.Children
		writer.Write([]byte("</"))
		writer.Write([]byte(n.Name.String()))
		writer.Write([]byte(">"))
	}
	return nil
}
