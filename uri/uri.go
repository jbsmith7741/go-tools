package uri

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/jbsmith7741/go-tools/appenderr"
)

var (
	// Separator used for slices
	Separator = ","

	// supported struct tags
	uriTag     = "uri"
	defaultTag = "default"

	// supported tag values
	scheme    = "scheme"
	host      = "host"
	path      = "path"
	authority = "authority" // scheme://host
	origin    = "origin"    // scheme://host/path
)

// Unmarshal copies a standard parsable uri to a predefined struct
// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
// scheme:opaque[?query][#fragment]
func Unmarshal(uri string, v interface{}) error {
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
	errs := appenderr.New()
	for i := 0; i < vStruct.NumField(); i++ {
		field := vStruct.Field(i)

		if !field.CanSet() { // skip private variables
			continue
		}

		name := vStruct.Type().Field(i).Name
		tag := vStruct.Type().Field(i).Tag.Get(uriTag)
		if tag != "" {
			name = tag
			tag = strings.ToLower(tag)
		}

		// check default values
		def := vStruct.Type().Field(i).Tag.Get(defaultTag)
		if def != "" {
			if err := SetField(field, def); err != nil {
				errs.Add(fmt.Errorf("default value %s can not be set to %s (%s)", def, name, field.Type()))
			}
		}

		data := values.Get(name)
		switch tag {
		case scheme:
			data = u.Scheme
		case host:
			data = u.Host
		case path:
			data = u.Path
		case origin:
			data = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
			if u.Scheme == "" && u.Host == "" {
				data = u.Path
			}
		case authority:
			data = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		default:
			if len(values[name]) == 0 {
				continue
			}
		}

		if field.Kind() == reflect.Slice {
			data = strings.Join(values[name], Separator)
		}

		if err := SetField(field, data); err != nil {
			errs.Add(fmt.Errorf("%s can not be set to %s (%s)", data, name, field.Type()))
		}
	}
	return errs.ErrOrNil()
}

func Marshal(v interface{}) (s string) {
	var u url.URL
	uVal := url.Values{}
	var vStruct reflect.Value
	if reflect.TypeOf(vStruct).Kind() == reflect.Ptr {
		vStruct = reflect.ValueOf(v).Elem()
	} else {
		vStruct = reflect.ValueOf(v)
	}

	for i := 0; i < vStruct.NumField(); i++ {
		field := vStruct.Field(i)

		var name string
		tag := vStruct.Type().Field(i).Tag.Get(uriTag)

		fs := GetFieldString(field)
		switch tag {
		case scheme:
			u.Scheme = fs
			continue
		case host:
			u.Host = fs
			continue
		case path:
			u.Path = fs
			continue
		case origin:
		case authority:
		case "":
			name = vStruct.Type().Field(i).Name
		default:
			name = tag
		}
		def := vStruct.Type().Field(i).Tag.Get(defaultTag)
		//fmt.Printf("%s|tag:%q,default:%q|%s|%v\n", name, tag, def, fs, field.Interface())
		// skip default fields
		if def == "" && isZero(field) {
			continue
		} else if fs == def {
			continue
		}

		if field.Kind() == reflect.Slice {
			for _, v := range strings.Split(fs, ",") {
				uVal.Add(name, v)
			}
		} else {
			uVal.Add(name, fs)
		}
	}

	// Note: url values are sorted by string value as they are encoded
	u.RawQuery = uVal.Encode()

	return u.String()
}

// SetField converts the string s to the type of value and sets the value if possible.
// Pointers and slices are recursively dealt with by deferencing the pointer
// or creating a generic slice of type value.
// All structs and alias' that implement encoding.TextUnmarshaler are suppported
func SetField(value reflect.Value, s string) error {
	if isAlias(value) {
		v := reflect.New(value.Type())
		if implementsUnmarshaler(v) {
			err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
			if err != nil {
				return err
			}
			value.Set(v.Elem())
			return nil
		}
	}
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
		SetField(z.Elem(), s)
		value.Set(z)

	case reflect.Slice:
		// create a generate slice and recursively assign the elements
		baseType := reflect.TypeOf(value.Interface()).Elem()
		data := strings.Split(s, Separator)
		slice := reflect.MakeSlice(value.Type(), 0, len(data))
		for _, v := range data {
			baseValue := reflect.New(baseType).Elem()
			SetField(baseValue, v)
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

func GetFieldString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.Interface().(string)
	case reflect.Bool:
		if value.Interface().(bool) == true {
			return "true"
		} else {
			return "false"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", value.Interface())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", value.Interface())
	case reflect.Ptr:
		if value.IsNil() {
			return "nil"
		}
		return GetFieldString(value.Elem())
	case reflect.Slice:
		var s string
		for i := 0; i < value.Len(); i++ {
			s += GetFieldString(value.Index(i)) + ","
		}
		return strings.TrimRight(s, ",")
	case reflect.Struct:
		s, _ := tryMarshal(value)
		return s
	default:
		return ""
	}
}

func isAlias(v reflect.Value) bool {
	if v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr {
		return false
	}
	s := fmt.Sprint(v.Type())
	return strings.Contains(s, ".")
}

func implementsUnmarshaler(v reflect.Value) bool {
	return v.Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem())
}

func tryMarshal(v reflect.Value) (string, error) {
	// does it implement TextMarshaler?
	if v.Type().Implements(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()) {
		b, err := v.Interface().(encoding.TextMarshaler).MarshalText()
		return string(b), err
	} else if v.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
		return v.Interface().(fmt.Stringer).String(), nil
	}
	return "", nil
}

func isZero(v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
