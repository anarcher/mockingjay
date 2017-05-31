package xml

import (
	"fmt"
	"reflect"
	"time"

	"github.com/beevik/etree"
	"github.com/fatih/structs"
)

const (
	sliceElementName = "member"
	locationName     = "locationName"
	locationNameList = "locationListName"
)

func NewElement(pe *etree.Element, v interface{}, names ...string) {

	//TODO(anarcher) Need refactoring
	if ts, ok := v.(*time.Time); ok {
		tsStr := ts.Format(time.RFC3339)
		ValueElement(pe, tsStr)
		return
	}

	switch Kind(v) {
	case reflect.Struct:
		StructElement(pe, v, names...)
	case reflect.String:
		ValueElement(pe, v)
	case reflect.Float64:
		ValueElement(pe, v)
	case reflect.Int64:
		ValueElement(pe, v)
	case reflect.Int:
		ValueElement(pe, v)
	case reflect.Bool:
		ValueElement(pe, v)
	case reflect.Slice:
		SliceElement(pe, v, names...)
	}
}

func StructElement(pe *etree.Element, v interface{}, names ...string) {
	s := structs.New(v)
	if structs.IsZero(s) {
		return
	}

	name := ElementName(s.Name(), names)
	pe = pe.CreateElement(name)

	for _, f := range s.Fields() {
		if !f.IsExported() || f.Name() == "_" || f.IsZero() {
			continue
		}

		name := FieldName(f, locationName)
		e := pe.CreateElement(name)
		locNameList := f.Tag(locationNameList)
		NewElement(e, f.Value(), locNameList)
	}
}

func ValueElement(e *etree.Element, v interface{}) {
	vv := Value(v)
	e.SetText(fmt.Sprintf("%v", vv))
}

func SliceElement(pe *etree.Element, v interface{}, names ...string) {
	name := ElementName(sliceElementName, names)

	s := Value(v)
	for i := 0; i < s.Len(); i++ {
		if !s.Index(i).IsValid() {
			continue
		}
		item := s.Index(i).Interface()
		//TODO(anarcher) It's not good. ValueElement actually is not Element
		if structs.IsStruct(item) {
			NewElement(pe, item, name)
		} else {
			e := pe.CreateElement(name)
			NewElement(e, item, name)
		}
	}

}
