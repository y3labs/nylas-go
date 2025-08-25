package nylas

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// EncodeQuery converts a struct or map into url.Values.
// - nil pointers are omitted
// - slices repeat the same key (k=v1&k=v2)
func EncodeQuery(params any) url.Values {
	v := url.Values{}
	if params == nil {
		return v
	}

	rv := reflect.ValueOf(params)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return v
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return v
	}
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		key := tagKey(sf)
		if key == "" || key == "-" {
			continue
		}

		fv := rv.Field(i)
		if !fv.IsValid() {
			continue
		}

		// Unwrap pointers; skip nil
		for fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				fv = reflect.Value{}
				break
			}
			fv = fv.Elem()
		}
		if !fv.IsValid() {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			if s := fv.String(); s != "" {
				v.Set(key, s)
			}
		case reflect.Bool:
			v.Set(key, strconv.FormatBool(fv.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.Set(key, strconv.FormatInt(fv.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v.Set(key, strconv.FormatUint(fv.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			v.Set(key, strconv.FormatFloat(fv.Float(), 'f', -1, 64))
		case reflect.Slice:
			n := fv.Len()
			if n == 0 {
				continue
			}
			elemKind := fv.Type().Elem().Kind()
			switch elemKind {
			case reflect.String:
				for j := 0; j < n; j++ {
					// Works for named string types too
					v.Add(key, fv.Index(j).String())
				}
			case reflect.Bool:
				for j := 0; j < n; j++ {
					v.Add(key, strconv.FormatBool(fv.Index(j).Bool()))
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				for j := 0; j < n; j++ {
					v.Add(key, strconv.FormatInt(fv.Index(j).Int(), 10))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				for j := 0; j < n; j++ {
					v.Add(key, strconv.FormatUint(fv.Index(j).Uint(), 10))
				}
			case reflect.Float32, reflect.Float64:
				for j := 0; j < n; j++ {
					v.Add(key, strconv.FormatFloat(fv.Index(j).Float(), 'f', -1, 64))
				}
			default:
				// Last resort: fmt.Sprint each element
				for j := 0; j < n; j++ {
					v.Add(key, fmt.Sprint(fv.Index(j).Interface()))
				}
			}
		default:
			// Fallback for custom named types; fmt.Sprint handles string aliases nicely
			s := fmt.Sprint(fv.Interface())
			if s != "" && s != "<nil>" {
				v.Set(key, s)
			}
		}
	}
	return v
}

func tagKey(sf reflect.StructField) string {
	if raw, ok := sf.Tag.Lookup("url"); ok && raw != "" {
		return strings.Split(raw, ",")[0]
	}
	if raw, ok := sf.Tag.Lookup("query"); ok && raw != "" {
		return strings.Split(raw, ",")[0]
	}
	return ""
}
