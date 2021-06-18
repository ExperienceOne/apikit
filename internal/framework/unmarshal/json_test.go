package unmarshal_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ExperienceOne/apikit/internal/framework/unmarshal"
)

type Constant string

const testConstant Constant = "test123"

const (
	nullJSON            = "null"
	testInt             = 123
	testFloatJSON       = `123`
	testFloatMapJSON    = `{ "Number": ` + testFloatJSON + `}`
	testString          = "test123"
	testStringJSON      = `"` + testString + `"`
	testSliceJSON       = `[ ` + testStringJSON + `]`
	testKey             = "key"
	testKeyJSON         = `"` + testKey + `"`
	testMapJSON         = `{` + testKeyJSON + `: ` + testStringJSON + `}`
	testSliceObjectJSON = `[` + testMapJSON + `]`
	testNullMap         = `{}`
	testMap             = "map"
	testMapKeyJSON      = `"` + testMap + `"`
	testSliceMapJSON    = `[` + `{` + testMapKeyJSON + `: {` + testKeyJSON + `: ` + testStringJSON + `}` + `}` + `]`
	testArray           = "array"
	testArrayKeyJSON    = `"` + testArray + `"`
	testSliceSliceJSON  = `[` + `{` + testArrayKeyJSON + `: ` + testSliceJSON + `}` + `]`
)

func TestJSON_constant(t *testing.T) {

	var s Constant
	err := unmarshal.JSON(strings.NewReader(testStringJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if s != testConstant {
		t.Errorf("expected '%s', got '%s'", testConstant, s)
	}

}

func TestJSON_object_constant(t *testing.T) {

	var s struct {
		Key Constant `json:"key"`
	}

	err := unmarshal.JSON(strings.NewReader(testMapJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if s.Key != testConstant {
		t.Errorf("expected '%s', got '%s'", testConstant, s.Key)
	}

}

func TestJSON_int(t *testing.T) {

	var s map[string]int

	err := unmarshal.JSON(strings.NewReader(testFloatMapJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if s["Number"] != testInt {
		t.Errorf("expected '%d', got '%d'", testInt, s["Number"])
	}
}

func TestJSON_float(t *testing.T) {

	var s map[string]float32

	err := unmarshal.JSON(strings.NewReader(testFloatMapJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if s["Number"] != testInt {
		t.Errorf("expected '%d', got '%f'", testInt, s["Number"])
	}
}

func TestJSON_string_required(t *testing.T) {

	var s string
	err := unmarshal.JSON(strings.NewReader(testStringJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if s != testString {
		t.Errorf("expected '%s', got '%s'", testString, s)
	}
}

func TestJSON_string_not_required(t *testing.T) {

	var s *string
	err := unmarshal.JSON(strings.NewReader(testStringJSON), &s, false)
	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("unexpected nil string")
	}

	if *s != testString {
		t.Errorf("expected '%s', got '%s'", testString, *s)
	}
}

func TestJSON_string_null(t *testing.T) {

	var s *string
	err := unmarshal.JSON(strings.NewReader(nullJSON), &s, false)
	if err != nil {
		t.Fatal(err)
	}

	if s != nil {
		t.Error("expected nil pointer")
	}
}

func TestJSON_string_null_but_required(t *testing.T) {

	var s string
	err := unmarshal.JSON(strings.NewReader(nullJSON), &s, true)
	if err == nil {
		t.Error("expected null error")
	} else if !errors.Is(err, unmarshal.NullError) {
		t.Error(err)
	}
}

func TestJSON_slice_required(t *testing.T) {

	var s []string
	err := unmarshal.JSON(strings.NewReader(testSliceJSON), &s, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(s) != 1 {
		t.Fatal("expected slice of len 1")
	}

	if s[0] != testString {
		t.Errorf("expected '%s', got '%s'", testString, s[0])
	}
}

func TestJSON_slice_not_required(t *testing.T) {

	var s []**string
	err := unmarshal.JSON(strings.NewReader(testSliceJSON), &s, false)
	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("unexpected nil slice")
	}

	if len(s) != 1 {
		t.Fatal("expected slice of len 1")
	}

	if s[0] == nil {
		t.Fatal("unexpected nil slice item")
	}

	if **(s[0]) != testString {
		t.Errorf("expected '%s', got '%s'", testString, **(s[0]))
	}
}

func TestJSON_slice_null(t *testing.T) {

	var s []string
	err := unmarshal.JSON(strings.NewReader(nullJSON), &s, false)
	if err != nil {
		t.Fatal(err)
	}

	if s != nil {
		t.Error("expected nil pointer")
	}
}

func TestJSON_slice_null_but_required(t *testing.T) {

	var s []string
	err := unmarshal.JSON(strings.NewReader(nullJSON), &s, true)
	if err == nil {
		t.Error("expected null error")
	} else if !errors.Is(err, unmarshal.NullError) {
		t.Error(err)
	}
}

func TestJSON_map_required(t *testing.T) {

	var s map[string]string
	err := unmarshal.JSON(strings.NewReader(testMapJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(s) != 1 {
		t.Fatal("expected map of len 1")
	}

	if s[testKey] != testString {
		t.Errorf("expected '%s', got '%s'", testString, s[testKey])
	}
}

func TestJSON_map_not_required(t *testing.T) {

	var s map[string]*string
	err := unmarshal.JSON(strings.NewReader(testMapJSON), &s, false)
	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("unexpected nil map")
	}

	if len(s) != 1 {
		t.Fatal("expected map of len 1")
	}

	if s[testKey] == nil {
		t.Fatal("unexpected nil map value")
	}

	if *(s[testKey]) != testString {
		t.Errorf("expected '%s', got '%s'", testString, *(s[testKey]))
	}
}

func TestJSON_map_null(t *testing.T) {

	var s map[string]string
	err := unmarshal.JSON(strings.NewReader(nullJSON), &s, false)
	if err != nil {
		t.Fatal(err)
	}

	if s != nil {
		t.Error("expected nil pointer")
	}
}

func TestJSON_map_null_but_required(t *testing.T) {

	var s map[string]string
	err := unmarshal.JSON(strings.NewReader(nullJSON), &s, true)
	if err == nil {
		t.Error("expected null error")
	} else if !errors.Is(err, unmarshal.NullError) {
		t.Error(err)
	}
}

func TestJSON_object_required(t *testing.T) {

	var s struct {
		Key string `json:"key"`
	}

	err := unmarshal.JSON(strings.NewReader(testMapJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if s.Key != testString {
		t.Errorf("expected '%s', got '%s'", testString, s.Key)
	}

}

func TestJSON_object_not_required(t *testing.T) {

	var s *struct {
		Key **string `json:"key"`
	}

	err := unmarshal.JSON(strings.NewReader(testMapJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("unexpected nil pointer")
	}

	if s.Key == nil || *s.Key == nil {
		t.Fatal("unexpected nil field")
	}

	if **s.Key != testString {
		t.Errorf("expected '%s', got '%s'", testString, **s.Key)
	}
}

func TestJSON_object_value_required_but_null(t *testing.T) {

	var s struct {
		Key string `json:"key,required"`
	}

	err := unmarshal.JSON(strings.NewReader(testNullMap), &s, true)
	if err == nil {
		t.Error("expected null error")
	} else if !errors.Is(err, unmarshal.NullError) {
		t.Error(err)
	}
}

func TestJSON_slice_object_required(t *testing.T) {

	var s []struct {
		Key string `json:"key"`
	}

	err := unmarshal.JSON(strings.NewReader(testSliceObjectJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(s) != 1 {
		t.Fatal("expected slice of len 1")
	}

	if s[0].Key != testString {
		t.Errorf("expected '%s', got '%s'", testString, s[0].Key)
	}

}

func TestJSON_slice_map_required(t *testing.T) {

	var s []struct {
		Map map[string]string `json:"map"`
	}

	err := unmarshal.JSON(strings.NewReader(testSliceMapJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(s) != 1 {
		t.Fatal("expected slice of len 1")
	}

	if len(s[0].Map) == 0 {
		t.Fatal("expected slice item map with len 1")
	}

	if s[0].Map[testKey] != testString {
		t.Errorf("expected '%s', got '%s'", testString, s[0].Map[testKey])
	}
}

func TestJSON_slice_slice_required(t *testing.T) {

	var s []struct {
		Array []string `json:"array"`
	}

	err := unmarshal.JSON(strings.NewReader(testSliceSliceJSON), &s, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(s) != 1 {
		t.Fatal("expected slice of len 1")
	}

	if len(s[0].Array) == 0 {
		t.Fatal("expected slice item slice with len 1")
	}

	if s[0].Array[0] != testString {
		t.Errorf("expected '%s', got '%s'", testString, s[0].Array[0])
	}
}

func TestJSON_panic(t *testing.T) {

	t.Log(unmarshal.JSON(nil, nil, false))
}
