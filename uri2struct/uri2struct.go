package uri2struct

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
// scheme:opaque[?query][#fragment]
func Convert(v interface{}, uri string) error {
	//verify that v is a pointer

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	values := u.Query()

	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("")
	}
	vStruct := reflect.ValueOf(v).Elem()
	for i := 0; i < vStruct.NumField(); i++ {
		field := vStruct.Field(i)
		name := field.Type().Name()
		s := values.Get(name)
		if s == "" { // skip fields not found
			continue
		}
		if err := setField(field, values.Get(name)); err != nil {
			return err
		}
	}
	return nil
}

func setField(value reflect.Value, s string) error {
	switch value.Kind() {
	case reflect.String:
		value.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return err
		}
		value.SetInt(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 0)
		if err != nil {
			return err
		}
		value.SetFloat(f)
	default:
		return fmt.Errorf("Unsupported type %v", value.Kind())
	}
	return nil
}
