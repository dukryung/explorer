package codec

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

var (
	cdc map[string]interface{}
)

func init() {
	cdc = make(map[string]interface{})
}

func RegisterCodec(path string, params interface{}) {
	cdc[fmt.Sprintf("/api%s", path)] = params
}

func Decode(p string, val url.Values) []interface{} {
	var values []interface{}
	if codec, ok := cdc[p]; ok {
		e := reflect.TypeOf(codec)

		for k, v := range val {
			for i := 0; i < e.NumField(); i++ {
				tag, ok := e.Field(i).Tag.Lookup("field")
				if !ok {
					continue
				}

				if tag != k {
					continue
				}

				o, ok := e.Field(i).Tag.Lookup("order")
				if !ok {
					continue
				}

				order, err := strconv.ParseInt(o, 10, 0)
				if err != nil {
					continue
				}

				value := parseValue(v[0], e.Field(i).Type)
				if value != nil {
					values = insertValue(values, int(order), value)
				}
			}
		}
	}

	return values
}

func insertValue(slice []interface{}, index int, value interface{}) []interface{} {
	if len(slice) <= index {
		return append(slice, value)
	}
	return append(slice[:index], append([]interface{}{value}, slice[index:]...)...)
}

func parseValue(param string, rtype reflect.Type) interface{} {
	value := reflect.New(rtype).Elem()

	switch value.Kind() {
	case reflect.Int32, reflect.Int8, reflect.Int16, reflect.Int64,
		reflect.Uint32, reflect.Uint8, reflect.Uint16, reflect.Uint64:
		ui, err := strconv.ParseInt(param, 10, 0)
		if err != nil {
			return nil
		}
		return ui
	case reflect.Float64, reflect.Float32:
		f, err := strconv.ParseFloat(param, 0)
		if err != nil {
			return nil
		}
		return f
	case reflect.String:
		return param
	case reflect.Slice:
		return nil
	}
	return nil
}
