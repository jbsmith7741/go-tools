package uri2struct

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var (
	// Seperator used for slices
	Seperator = ","

	// supported struct tags
	uriTag = "uri"

	// supported tag values
	scheme = "scheme"
	host   = "host"
	path   = "path"
	origin = "origin" // schema://host/path
)

// Convert copies a standard parsable uri to a predefined struct
// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
// scheme:opaque[?query][#fragment]
func Convert(v interface{}, uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	values := u.Query()

	//verify that v is a pointer
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("struct must be a pointer")
	}
	vStruct := reflect.ValueOf(v).Elem()
	for i := 0; i < vStruct.NumField(); i++ {
		field := vStruct.Field(i)

		name := vStruct.Type().Field(i).Name
		tag, found := vStruct.Type().Field(i).Tag.Lookup(uriTag)
		if found {
			name = tag
			tag = strings.ToLower(tag)
		}

		data := values.Get(name)
		if tag == scheme {
			data = u.Scheme
		} else if tag == host {
			data = u.Host
		} else if tag == path {
			data = u.Path
		} else if tag == origin {
			data = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
		} else if len(values[name]) == 0 {
			continue
		}

		if field.Kind() == reflect.Slice {
			// TODO: should this be default behavior?
			data = strings.Join(values[name], Seperator)
		}

		if err := setField(field, data); err != nil {
			return err
		}
	}
	return nil
}

/* func checkStructTag(v interface) string {

} */

func setField(value reflect.Value, s string) error {
	switch value.Kind() {
	case reflect.String:
		value.SetString(s)
	case reflect.Bool:
		b := strings.ToLower(s) == "true" || s == ""
		value.SetBool(b)
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

	case reflect.Ptr:
		// create non pointer type and recursively assign
		z := reflect.New(value.Type().Elem())
		setField(z.Elem(), s)
		value.Set(z)

	case reflect.Slice:
		// create a generate slice and recursively assign the elements
		baseType := reflect.TypeOf(value.Interface()).Elem()
		data := strings.Split(s, Seperator)
		slice := reflect.MakeSlice(value.Type(), 0, len(data))
		for _, v := range data {
			baseValue := reflect.New(baseType).Elem()
			setField(baseValue, v)
			slice = reflect.Append(slice, baseValue)
		}
		value.Set(slice)

	case reflect.Struct:
		v := reflect.New(value.Type())
		if implementsUnmarshaler(v) {
			err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
			if err != nil {
				return err
			}
		}
		value.Set(v.Elem())

	default:
		return fmt.Errorf("Unsupported type %v", value.Kind())
	}
	return nil
}

func implementsUnmarshaler(v reflect.Value) bool {
	return v.Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem())
}
