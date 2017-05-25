package xml

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/beevik/etree"
)

func buildXML(name string, v interface{}) (string, error) {
	doc := etree.NewDocument()
	e := doc.CreateElement(name)
	NewElement(e, v)
	xml, err := doc.WriteToString()
	return xml, err
}

func TestNewElement(t *testing.T) {
	type (
		Test struct {
			Name string
			Ok   bool
		}
		Test2 struct {
			Name *string
		}
	)

	var tests = []struct {
		doc      string
		value    interface{}
		expected string
	}{
		{
			"TestCase",
			&Test{
				Name: "test1",
			},
			"<TestCase><Test><Name>test1</Name></Test></TestCase>",
		},
		{
			"T",
			&Test2{
				Name: aws.String("t"),
			},
			"<T><Test2><Name>t</Name></Test2></T>",
		},
	}

	for i, test := range tests {

		xml, err := buildXML(test.doc, test.value)
		if err != nil {
			t.Error(err)
		}

		if xml != test.expected {
			t.Errorf("%d: have: %v want:%v", i, xml, test.expected)
		}
	}
}

func TestSliceElement(t *testing.T) {
	doc := etree.NewDocument()
	e := doc.CreateElement("E")
	v := []string{"a", "b", "c"}
	SliceElement(e, v, "I")

	xml, err := doc.WriteToString()
	if err != nil {
		t.Error(err)
	}

	expected := "<E><I>a</I><I>b</I><I>c</I></E>"
	if xml != expected {
		t.Errorf("have: %v want: %v", xml, expected)
	}
}
