package parameter_test

import (
	"reflect"
	"testing"

	"errors"
	"strings"

	"github.com/ExperienceOne/apikit/internal/framework/parameter"
)

func TestToString(t *testing.T) {

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "test", "test"},
		{"[]string", []string{"test1", "test2"}, "test1,test2"},
		{"true", true, "true"},
		{"false", false, "false"},
		{"[]bool", []bool{true, false}, "true,false"},
		{"int", int(123), "123"},
		{"int8", int8(123), "123"},
		{"int16", int16(123), "123"},
		{"int32", int32(123), "123"},
		{"int64", int64(123), "123"},
		{"uint", uint(123), "123"},
		{"uint8", uint8(123), "123"},
		{"uint16", uint16(123), "123"},
		{"uint32", uint32(123), "123"},
		{"uint64", uint64(123), "123"},
		{"[]int", []int{int(12), int(34)}, "12,34"},
		{"[]int8", []int8{int8(12), int8(34)}, "12,34"},
		{"[]int16", []int16{int16(12), int16(34)}, "12,34"},
		{"[]int32", []int32{int32(12), int32(34)}, "12,34"},
		{"[]int64", []int64{int64(12), int64(34)}, "12,34"},
		{"[]uint", []uint{uint(12), uint(34)}, "12,34"},
		{"[]uint8", []uint8{uint8(12), uint8(34)}, "12,34"},
		{"[]uint16", []uint16{uint16(12), uint16(34)}, "12,34"},
		{"[]uint32", []uint32{uint32(12), uint32(34)}, "12,34"},
		{"[]uint64", []uint64{uint64(12), uint64(34)}, "12,34"},
		{"float32", float32(123.45), "123.45"},
		{"[]float32", []float32{float32(123.45), float32(67.89)}, "123.45,67.89"},
		{"float64", float64(123.450), "123.45"},
		{"[]float64", []float64{float64(123.45), float64(67.89)}, "123.45,67.89"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := parameter.ToString(test.input)
			if result != test.expected {
				t.Fatalf("expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestFromString(t *testing.T) {

	type customInt int
	emptyString := ""

	tests := []struct {
		name     string
		input    string
		typ      reflect.Type
		expected interface{}
	}{
		{"[]customInt", "1,1,3,5,7", reflect.TypeOf([]customInt{}), []customInt{1, 1, 3, 5, 7}},
		{"customInt", "1", reflect.TypeOf(customInt(1)), customInt(1)},
		{"string", "test", reflect.TypeOf(emptyString), "test"},
		{"*string", "test", reflect.TypeOf(&emptyString), "test"},
		{"[]string", "test1,test2", reflect.TypeOf([]string{}), []string{"test1", "test2"}},
		{"true", "True", reflect.TypeOf(true), true},
		{"false", "False", reflect.TypeOf(false), false},
		{"bool", "true,false", reflect.TypeOf([]bool{}), []bool{true, false}},
		{"int8", "123", reflect.TypeOf(int8(0)), int8(123)},
		{"int16", "123", reflect.TypeOf(int16(0)), int16(123)},
		{"int32", "123", reflect.TypeOf(int32(0)), int32(123)},
		{"int64", "123", reflect.TypeOf(int64(0)), int64(123)},
		{"uint", "123", reflect.TypeOf(uint(0)), uint(123)},
		{"uint8", "123", reflect.TypeOf(uint8(0)), uint8(123)},
		{"uint16", "123", reflect.TypeOf(uint16(0)), uint16(123)},
		{"uint32", "123", reflect.TypeOf(uint32(0)), uint32(123)},
		{"uint64", "123", reflect.TypeOf(uint64(0)), uint64(123)},
		{"[]int", "123,213", reflect.TypeOf([]int{}), []int{int(123), int(213)}},
		{"[]int8", "123,56", reflect.TypeOf([]int8{}), []int8{int8(123), int8(56)}},
		{"[]int16", "123,213", reflect.TypeOf([]int16{}), []int16{int16(123), int16(213)}},
		{"[]int32", "123,213", reflect.TypeOf([]int32{}), []int32{int32(123), int32(213)}},
		{"[]int64", "123,213", reflect.TypeOf([]int64{}), []int64{int64(123), int64(213)}},
		{"[]uint", "123,213", reflect.TypeOf([]uint{}), []uint{uint(123), uint(213)}},
		{"[]uint8", "123,213", reflect.TypeOf([]uint8{}), []uint8{uint8(123), uint8(213)}},
		{"[]uint16", "123,213", reflect.TypeOf([]uint16{}), []uint16{uint16(123), uint16(213)}},
		{"[]uint32", "123,213", reflect.TypeOf([]uint32{}), []uint32{uint32(123), uint32(213)}},
		{"[]uint64", "123,213", reflect.TypeOf([]uint64{}), []uint64{uint64(123), uint64(213)}},
		{"float32", "123.45", reflect.TypeOf(float32(0)), float32(123.45)},
		{"float64", "123.45", reflect.TypeOf(float64(0)), float64(123.45)},
		{"[]float32", "123.45,213.45", reflect.TypeOf([]float32{}), []float32{float32(123.45), float32(213.45)}},
		{"[]float64", "123.45,213.45", reflect.TypeOf([]float64{}), []float64{float64(123.45), float64(213.45)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variable := reflect.New(test.typ).Interface()
			err := parameter.FromString(test.input, variable)
			if err != nil {
				t.Fatal(err)
			}
			if reflect.ValueOf(variable).Elem().Kind() == reflect.Ptr {
				variable = reflect.ValueOf(variable).Elem().Interface()
			}
			value := reflect.ValueOf(variable).Elem().Interface()
			if !reflect.DeepEqual(value, test.expected) {
				t.Fatalf("expected '%v', got '%v'", test.expected, value)
			}
		})
	}
}

func TestFromStringFailed(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		typ      reflect.Type
		expected error
	}{
		{"[9]string", "1,2,3", reflect.TypeOf([9]string{}), errors.New("unsupported kind: array")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variable := reflect.New(test.typ).Interface()
			err := parameter.FromString(test.input, variable)
			if err == nil {
				t.Fatal("err is nil")
			}

			if err == test.expected {
				t.Fatalf("expected '%v', got '%v'", test.expected, err)
			}
		})
	}
}

func TestFromStringIsNotAPointer(t *testing.T) {

	err := parameter.FromString("1", nil)
	if err == nil {
		t.Fatal("err is nil")
	}

	if !errors.Is(err, parameter.ErrParamIsNotPointer) {
		t.Fatalf("expected '%v', got '%v'", parameter.ErrParamIsNotPointer, err)
	}
}

func TestFormStringNil(t *testing.T) {
	var nilInt *int
	err := parameter.FromString("1", nilInt)
	if err == nil {
		t.Fatal("error is nil")
	}

	if !errors.Is(err, parameter.ErrParamIsNil) {
		t.Fatalf("expected '%v', got '%v'", parameter.ErrParamIsNil, err)
	}
}

func TestFromStringOverflow(t *testing.T) {
	var uintVal uint8
	var intVal int8
	tests := []struct {
		name  string
		value interface{}
		raw   string
	}{
		{
			name:  "unit8",
			value: &uintVal,
			raw:   "1000",
		},

		{
			name:  "int8",
			value: &intVal,
			raw:   "1000",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := parameter.FromString(test.raw, test.value)
			if err == nil {
				t.Fatal("err is nil")
			}
			if !strings.Contains(err.Error(), "type overflow:") {
				t.Errorf("didn't triggered a overflow (message: %v)", err)
			}
		})
	}
}
