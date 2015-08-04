package xtraml

import (
	"encoding/xml"
	"errors"
	"io"
)

func first(reader io.Reader) (*StartToken, error) {
	d := xml.NewDecoder(reader)
	base := baseToken{nil, d}
	for true {
		token := base.Next()
		if token == nil {
			return nil, io.EOF
		}
		start, ok := token.(StartToken)
		if ok {
			return &start, nil
		}
		base = base.copy()
	}
	panic(errors.New("Should not reach!"))
}

type Token interface {
	Next() Token
}

type baseToken struct {
	next Token
	d    *xml.Decoder
}

func (b baseToken) Next() Token {
	if b.next == nil {
		next, err := b.nextToken()
		if err != nil {
			panic(err)
		}
		b.next = next
	}
	return b.next
}

func (b baseToken) nextToken() (Token, error) {
	token, err := b.d.Token()
	if err != nil {
		if err == io.EOF {
			return EOF{}, nil
		}
		return nil, err
	}
	switch token.(type) {
	case xml.StartElement:
		start, _ := token.(xml.StartElement)
		return fromStart(start, b.copy()), nil
	case xml.EndElement:
		end, _ := token.(xml.EndElement)
		return fromEnd(end, b.copy()), nil
	case xml.CharData:
		charData, _ := token.(xml.CharData)
		return fromCharData(charData, b.copy()), nil
	case xml.Comment:
		comment, _ := token.(xml.Comment)
		return fromComment(comment, b.copy()), nil
	case xml.ProcInst:
		procInst, _ := token.(xml.ProcInst)
		return fromProcInst(procInst, b.copy()), nil
	default:
		return nil, errors.New("Invalid Token Type")
	}
	return nil, nil
}

func (b baseToken) copy() baseToken {
	return baseToken{nil, b.d}
}

type StartToken struct {
	baseToken
	Name  Name
	Attrs AttrSet
	end   xml.EndElement
}

func fromStart(start xml.StartElement, baseToken baseToken) Token {
	return StartToken{baseToken, Name{start.Name}, fromAttrs(start.Attr), start.End()}
}

type EndToken struct {
	baseToken
	Name Name
	end  xml.EndElement
}

func fromEnd(end xml.EndElement, baseToken baseToken) Token {
	return EndToken{baseToken, Name{end.Name}, end}
}

type Name struct {
	xml.Name
}

func (n Name) String() string {
	return n.Space + ":" + n.Local
}

type AttrSet map[Name]string

func fromAttrs(attrs []xml.Attr) AttrSet {
	set := AttrSet{}
	for _, attr := range attrs {
		set[Name{attr.Name}] = attr.Value
	}
	return set
}

func (a AttrSet) String() string {
	out := ""
	for name, value := range a {
		out += " " + name.String() + "=\"" + value + "\""
	}
	return out
}

type Comment struct {
	baseToken
	Data string
}

func fromComment(comment xml.Comment, baseToken baseToken) Token {
	return Comment{baseToken, string(comment)}
}

type Chardata struct {
	baseToken
	Data string
}

func fromCharData(charData xml.CharData, baseToken baseToken) Token {
	return Chardata{baseToken, string(charData)}
}

type ProcInst struct {
	baseToken
	Target string
	Inst   string
}

func fromProcInst(procInst xml.ProcInst, baseToken baseToken) Token {
	return ProcInst{baseToken, procInst.Target, string(procInst.Inst)}
}

type EOF struct{}

func (e EOF) Next() Token {
	return nil
}
