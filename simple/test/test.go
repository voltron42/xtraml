package main

import (
	"../"
	"encoding/xml"
	"fmt"
)

func main() {
	bytes := []byte("<q text=\"q1\"><yes><a text=\"a1\"/></yes><no><q test=\"q1\"><yes><a text=\"a2\"/></yes><no><a text=\"a3\"/></no></q></no></q>")
	node := simple.XmlNode{}
	err := xml.Unmarshal(bytes, &node)
	if err != nil {
		panic(err)
	}
	fmt.Println(node.ToXml())

}
