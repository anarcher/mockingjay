package xml

import (
	"reflect"

	"github.com/fatih/structs"
)

func IsStruct(v interface{}) bool {
	return structs.IsStruct(v)
}

func IsString(v interface{}) bool {
	if Value(v).Kind() == reflect.String {
		return true
	}

	return false
}

func Kind(v interface{}) reflect.Kind {
	if IsStruct(v) {
		return reflect.Struct
	}
	vv := Value(v)
	return vv.Kind()
}

func Value(v interface{}) reflect.Value {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	return vv

}

func FieldName(f *structs.Field, tagName string) string {
	name := f.Name()
	if f.Tag(tagName) != "" {
		name = f.Tag(tagName)
	}
	return name
}

func ElementName(name string, names []string) string {
	if len(names) > 0 {
		if names[0] != "" {
			return names[0]
		}
	}
	return name
}
