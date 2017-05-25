package xml

import (
	"fmt"

	"github.com/beevik/etree"
)

func Response(name string, v interface{}, xmlns string) (string, error) {
	doc := etree.NewDocument()
	if xmlns != "" {
		doc.CreateAttr("xmlns", xmlns)
	}
	e := doc.CreateElement(fmt.Sprintf("%sResponse", name))
	NewElement(e, v, fmt.Sprintf("%sResult", name))
	xml, err := doc.WriteToString()
	return xml, err
}
