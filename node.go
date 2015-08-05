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
	writer.Write([]byte(n.Attrs.String()))
	if len(n.Children) < 1 {
		writer.Write([]byte("/>"))
	} else {
		writer.Write([]byte(">"))
		for _, child := range n.Children {
			child.WriteXML(writer)
		}
		writer.Write([]byte("</"))
		writer.Write([]byte(n.Name.String()))
		writer.Write([]byte(">"))
	}
	return nil
}

type CommentNode string

func (c CommentNode) WriteXML(writer io.Writer) error {
	writer.Write([]byte("<!--" + c + "-->"))
	return nil
}

type CharDataNode string

func (c CharDataNode) WriteXML(writer io.Writer) error {
	writer.Write([]byte(c))
	return nil
}

type ProcInstNode struct {
	Target string
	Inst   string
}

func (p ProcInstNode) WriteXML(writer io.Writer) error {
	writer.Write([]byte("<?" + p.Target + " " + p.Inst + "?>"))
	return nil
}
