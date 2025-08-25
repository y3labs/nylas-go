package nylas

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type sstr string

func (s sstr) String() string { return "S:" + string(s) }

type qStruct struct {
	A *int    `query:"a"`
	B *string `query:"b"`
	C []int   `query:"c"`
}

func TestEncodeQuery(t *testing.T) {
	a := 1
	b := "x"
	vals := EncodeQuery(&qStruct{A: &a, B: &b, C: []int{1, 2}})
	if vals.Get("a") != "1" || vals.Get("b") != "x" {
		t.Fatalf("unexpected query: %#v", vals)
	}
	cs := vals["c"]
	if len(cs) != 2 {
		t.Fatalf("want 2 c params, got %v", cs)
	}
}

func TestToString_Edges(t *testing.T) {
	var p *int
	if toString(nil) != "" {
		t.Fatal("nil")
	}
	if toString(true) != "true" {
		t.Fatal("bool")
	}
	if toString(123) != "123" {
		t.Fatal("int")
	}
	if toString(1.5) != "1.5" {
		t.Fatal("float")
	}
	if toString([]byte("x")) != "x" {
		t.Fatal("bytes")
	}
	if toString(sstr("ok")) != "S:ok" {
		t.Fatal("Stringer")
	}
	if toString(p) != "" {
		t.Fatal("nil ptr")
	}
	i := 9
	if toString(&i) != "9" {
		t.Fatal("ptr deref")
	}
	if toString("hi") != "hi" {
		t.Fatal("string")
	}
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	// Fast paths.
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	case []byte:
		return string(t)
	}

	// Deref pointers.
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ""
		}
		return toString(rv.Elem().Interface())
	}

	// Common scalar kinds (avoids surprises like "%!f(float64=...)" etc.)
	switch rv.Kind() {
	case reflect.String:
		return rv.String()
	case reflect.Bool:
		if rv.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64)
	default:
		// Fallback: fmt.Sprint
		return fmt.Sprint(v)
	}
}
