package encode

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// setField converts the string s to the type of value and sets the value if possible.
// Pointers and slices are recursively dealt with by following the pointer
// or creating a generic slice of type value.
//
// All structs and that implement encoding.TextUnmarshaler are supported
//
// Does not support array literals.
func SetField(value reflect.Value, s string, sField reflect.StructField) error {
	if isZero(value.Kind(), s) {
		return nil
	}
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
		switch strings.ToLower(s) {
		case "true":
			value.SetBool(true)
		case "false", "":
			value.SetBool(false)
		default:
			// the bool value should be explicit to tell user
			// something is amiss.
			return fmt.Errorf("cannot assign '%v' to bool type", s)
		}
	case reflect.Int64:
		// check for time.Duration int64 special case.
		// support int64 value and duration as a string
		if value.Type().String() == "time.Duration" {
			d, err := time.ParseDuration(s)
			if err != nil {
				i, e2 := strconv.ParseInt(s, 10, 64)
				if e2 != nil {
					return err
				}
				d = time.Duration(i)
			}

			value.SetInt(int64(d))
			return nil
		}
		// handle normal int64 with other ints
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		i, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return err
		}

		value.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, 0)
		if err != nil {
			return err
		}

		value.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 0)
		if err != nil {
			return err
		}
		value.SetFloat(f)
	case reflect.Ptr:
		if isZero(value.Type().Elem().Kind(), s) {
			return nil
		}
		// create non pointer type and recursively assign
		z := reflect.New(value.Type().Elem())
		err := SetField(z.Elem(), s, sField)
		if err != nil {
			return err
		}

		value.Set(z)
	case reflect.Slice:
		// create a slice and recursively assign the elements
		baseType := reflect.TypeOf(value.Interface()).Elem()
		s = strings.Trim(s, "[]") // trim brackets for bracket support.
		vals := strings.Split(s, ",")

		slice := reflect.MakeSlice(value.Type(), 0, len(vals))
		for _, v := range vals {
			// trim whitespace from each value to support comma-separated with spaces.
			v = strings.TrimSpace(v)
			v = strings.Trim(v, `"'`)

			// each item must be the correct type.
			baseValue := reflect.New(baseType).Elem()
			err := SetField(baseValue, v, sField)
			if err != nil {
				return err
			}
			slice = reflect.Append(slice, baseValue)
		}

		value.Set(slice)
	case reflect.Array:
		s = strings.Trim(s, "[]") // trim brackets for bracket support.
		vals := strings.Split(s, ",")
		if value.Len() != len(vals) {
			return fmt.Errorf("cannot set array of different lengths got %d want %d", value.Len(), len(vals))
		}
		for i, v := range vals {
			if err := SetField(value.Index(i), v, sField); err != nil {
				return err
			}
		}

	// structs as values are simply ignored. They don't map cleanly for environment variables.
	case reflect.Struct:
		v := reflect.New(value.Type())
		if value.Type().String() == "time.Time" {
			timeFmt := sField.Tag.Get(fmtTag)

			_, err := SetTime(value, s, timeFmt)
			return err
		}
		if implementsUnmarshaler(v) {
			err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
			if err != nil {
				return err
			}
		}
		value.Set(v.Elem())
		return nil
	default:
		return fmt.Errorf("unsupported type '%v'", value.Kind())
	}

	return nil
}

// isZero checks if the value s is the zero value of type t
func isZero(t reflect.Kind, s string) bool {
	switch t {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Float32, reflect.Float64:
		return s == "0"
	}
	return s == ""
}

// setTime expects value to be time.Time.
//
// tFmt can be any time package handy time format like "RFC3339Nano".
// Default format is time.RFC3339.
func SetTime(value reflect.Value, tv, timeFmt string) (string, error) {
	if timeFmt == "" {
		timeFmt = time.RFC3339 // default format
	}

	// check for standard time formats
	switch timeFmt {
	case "ANSIC":
		timeFmt = time.ANSIC
	case "UnixDate":
		timeFmt = time.UnixDate
	case "RubyDate":
		timeFmt = time.RubyDate
	case "RFC822":
		timeFmt = time.RFC822
	case "RFC822Z":
		timeFmt = time.RFC822Z
	case "RFC850":
		timeFmt = time.RFC850
	case "RFC1123":
		timeFmt = time.RFC1123
	case "RFC1123Z":
		timeFmt = time.RFC1123Z
	case "RFC3339":
		timeFmt = time.RFC3339
	case "RFC3339Nano":
		timeFmt = time.RFC3339Nano
	case "Kitchen":
		timeFmt = time.Kitchen
	case "Stamp":
		timeFmt = time.Stamp
	case "StampMilli":
		timeFmt = time.StampMilli
	case "StampMicro":
		timeFmt = time.StampMicro
	case "StampNano":
		timeFmt = time.StampNano
	}

	t, err := time.Parse(timeFmt, tv)
	if err != nil {
		return timeFmt, err
	}

	tStruct := reflect.ValueOf(t)
	value.Set(tStruct)

	return timeFmt, nil
}

func implementsUnmarshaler(v reflect.Value) bool {
	return v.Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem())
}

func isAlias(v reflect.Value) bool {
	if v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr {
		return false
	}
	s := fmt.Sprint(v.Type())
	return strings.Contains(s, ".")
}
