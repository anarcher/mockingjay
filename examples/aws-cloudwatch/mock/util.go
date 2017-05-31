package main

import (
	"fmt"
	"net/http"
)

func membersToMap(r *http.Request, name string) map[string]string {
	ret := make(map[string]string)

	for i := 1; i <= 10; i++ {
		keyname := fmt.Sprintf("%s.member.%d.Name", name, i)
		key := r.FormValue(keyname)
		if key != "" {
			valname := fmt.Sprintf("%s.member.%d.Value", name, i)
			value := r.FormValue(valname)
			ret[key] = value
		}
	}

	return ret
}
