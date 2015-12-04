package choice

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type ChoiceParser map[string]func(d *xml.Decoder, start xml.StartElement) (interface{}, error)

func (c ChoiceParser) Parse(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
	choice, ok := c[start.Name.Local]
	if !ok {
		return nil, errors.New(start.Name.Local + " is not a listed option.")
	}
	return choice(d, start)
}

func (c ChoiceParser) ParseList(d *xml.Decoder, start xml.StartElement, append func(item interface{}) error) error {
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
			err = append(item)
			if err != nil {
				return err
			}
		}
		token, err = d.Token()
	}
	return nil
}
