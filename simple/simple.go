package simple

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

type Node interface {
	ToXml() string
	AsXmlNode() (*XmlNode, error)
	Path(indices ...int) (Node, error)
}

type XmlNode struct {
	Name     string            `json:"name,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children []Node            `json:"children,omitempty"`
}

func (x *XmlNode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x.Name = getName(start.Name)
	x.Attrs = map[string]string{}
	for _, attr := range start.Attr {
		x.Attrs[getName(attr.Name)] = attr.Value
	}
	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}
		next, ok := token.(xml.StartElement)
		if ok {
			child := XmlNode{}
			err = child.UnmarshalXML(d, next)
			if err != nil {
				return err
			}
			x.Children = append(x.Children, child)
		} else {
			text, ok := token.(xml.CharData)
			if ok {
				x.Children = append(x.Children, TextNode(string([]byte(text))))
			}
		}
		token, err = d.Token()
	}
	return nil
}

func getName(name xml.Name) string {
	if len(name.Space) == 0 {
		return name.Local
	} else {
		return name.Space + ":" + name.Local
	}
}

func (x XmlNode) ToXml() string {
	attrtpl := " %v=\"%v\""
	attrs := []string{}
	for name, attr := range x.Attrs {
		attrs = append(attrs, fmt.Sprintf(attrtpl, name, attr))
	}
	children := []string{}
	for _, child := range x.Children {
		children = append(children, child.ToXml())
	}
	attrList := strings.Join(attrs, "")
	if len(children) > 0 {
		return fmt.Sprintf("<%v%v>%v</%v>", x.Name, attrList, strings.Join(children, ""), x.Name)
	} else {
		return fmt.Sprintf("<%v%v/>", x.Name, attrList)
	}
}

func (x XmlNode) AsXmlNode() (*XmlNode, error) {
	return &x, nil
}

func (x XmlNode) Path(indices ...int) (Node, error) {
	if len(indices) == 0 {
		return x, nil
	}
	if len(x.Children) <= indices[0] || indices[0] < 0 {
		return nil, errors.New("index out of bounds")
	}
	return x.Children[indices[0]].Path(indices[1:]...)
}

type TextNode string

func (t TextNode) ToXml() string {
	return string(t)
}

func (t TextNode) AsXmlNode() (*XmlNode, error) {
	return nil, errors.New("TextNode cannot be converted to an XmlNode")
}

func (t TextNode) Path(indices ...int) (Node, error) {
	if len(indices) > 0 {
		return nil, errors.New("text node cannot have a path")
	}
	return t, nil
}
