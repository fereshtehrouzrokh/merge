package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Merge(a, b map[string]interface{}) map[string]interface{} {
	if a == nil {
		a = make(map[string]interface{})
	}
	if b == nil || len(b) == 0 {
		return a
	}

	for k, bv := range b {
		av, aExists := a[k]
		switch {
		case aExists && isArray(av) && bv != nil && !isArray(bv):
			a[k] = append(av.([]interface{}), bv)
		case aExists && isArray(av) && bv == nil:
			a[k] = []interface{}{}
		case aExists && isArray(av) && isArray(bv):
			a[k] = bv
		case aExists && bv == nil:
			delete(a, k)
		case isMap(av) && isMap(bv):
			a[k] = Merge(av.(map[string]interface{}), bv.(map[string]interface{}))
		case bv != nil && isPrimitive(bv):
			a[k] = bv
		default:
			a[k] = bv
		}
	}
	return a
}
func isArray(v interface{}) bool {
	_, ok := v.([]interface{})
	return ok
}

func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}
func isPrimitive(v interface{}) bool {
	if v == nil {
		return false
	}
	k := reflect.TypeOf(v).Kind()
	return k >= reflect.Bool && k <= reflect.String
}

func main() {

	jsonA := []byte(`{
	  "first_name":"Bob",
	  "last_name":"Jones",
	  "email":"bob@gmail.com",
	  "address":{"line_1":"1234 Main St","line_2":"Apt 413","city":"Los Angeles","state":"CA","zip":"90048"},
	  "logins":[
	    {"date":"10/22/2012","ip":"192.168.0.1"},
	    {"date":"10/21/2012","ip":"192.168.0.1"}
	  ],
	  "photos":["IMG-1985.jpg","IMG-1987.jpg"]
	}`)

	jsonB := []byte(`{
	  "last_name":"Jones",
	  "active":true,
	  "address":{"line_1":"2143 South Main St","line_2":null},
	  "logins":{"date":"10/23/2012","ip":"192.168.0.1"},
	  "photos":null
	}`)


	var a, b map[string]interface{}
	if err := json.Unmarshal(jsonA, &a); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(jsonB, &b); err != nil {
		panic(err)
	}

	merged := Merge(a, b)


	out, _ := json.MarshalIndent(merged, "", "  ")
	fmt.Println(string(out))
}
