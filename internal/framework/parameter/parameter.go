package parameter

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// primitiveToString converts a given primitive value into an string
func primitiveToString(param reflect.Value) string {

	var value string

	switch param.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = fmt.Sprintf("%d", param.Interface())
	case reflect.Float64:
		value = strconv.FormatFloat(param.Interface().(float64), 'f', -1, 64)
	case reflect.Float32:
		value = strconv.FormatFloat(float64(param.Interface().(float32)), 'f', -1, 32)
	case reflect.String:
		value = fmt.Sprintf("%s", param.Interface())
	case reflect.Bool:
		value = fmt.Sprintf("%t", param.Interface())
	}

	return value
}



// sliceToString converts a given slice value into an string
func sliceToString(param reflect.Value) string {

	slice := make([]string, param.Len())
	for i := 0; i < param.Len(); i++ {
		slice[i] = ToString(param.Index(i).Interface())
	}
	return strings.Join(slice, ",")
}

// ToString converts a given value into an string
func ToString(param interface{}) string {

	paramReflected := reflect.ValueOf(param)

	for paramReflected.Kind() == reflect.Ptr {
		if paramReflected.IsNil() {
			return ""
		}
		paramReflected = paramReflected.Elem()
	}

	var value string
	if paramReflected.Kind() == reflect.Slice || paramReflected.Kind() == reflect.Array {
		value = sliceToString(paramReflected)
	} else {
		value = primitiveToString(paramReflected)
	}

	return value
}

// stringToPrimitive injects and converts an string value into an typed value
// if param is not a pointer then then return an error
// if the int to value too big then return an error
// if the uint value too big then return an error
func stringToPrimitive(s string, param reflect.Value) error {

	if param.Kind() != reflect.Ptr {
		return &ErrValueIsNotPointer{value: param}
	}

	var err error
	elm := param.Elem()

	switch elm.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := int64(0)
		if s != "" {
			val, err = strconv.ParseInt(s, 0, 64)
			if err != nil {
				return err
			}
			if elm.OverflowInt(val) {
				return &ErrTypeValueOverflow{value: s}
			}
		}
		elm.SetInt(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := uint64(0)
		if s != "" {
			val, err = strconv.ParseUint(s, 0, 64)
			if err != nil {
				return err
			}
			if elm.OverflowUint(val) {
				return &ErrTypeValueOverflow{value: s}
			}
		}
		elm.SetUint(val)

	case reflect.Float32:
		val := float64(0)
		if s != "" {
			val, err = strconv.ParseFloat(s, 32)
			if err != nil {
				return err
			}
		}
		elm.SetFloat(val)

	case reflect.Float64:
		val := float64(0)
		if s != "" {
			val, err = strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
		}
		elm.SetFloat(val)

	case reflect.String:
		elm.SetString(s)

	case reflect.Bool:
		val := false
		if s != "" {
			val, err = strconv.ParseBool(s)
			if err != nil {
				return err
			}
		}
		elm.SetBool(val)

	default:
		return &ErrUnsupportedPrimitiveType{value: param}
	}
	return nil
}

// stringToSlice injects and converts a string list value (separated by ",") into an typed slice
func stringToSlice(s string, param reflect.Value) error {

	values := strings.Split(s, ",")
	elemForInjection := reflect.New(param.Elem().Type().Elem())
	for _, value := range values {
		err := FromString(value, elemForInjection.Interface())
		if err != nil {
			return err
		}
		slice := reflect.Append(param.Elem(), elemForInjection.Elem())
		param.Elem().Set(slice)
	}

	return nil
}

// FromString injects and converts a string value into a spefific typed value
// if a panic occurred then return an error
// if the param is nil then return an error
func FromString(s string, param interface{}) (err error) {

	defer func() {
		if v := recover(); v != nil {
			stack := debug.Stack()
			err = &RecoverError{Err: v, Stack: stack}
		}
	}()

	for {
		paramReflected := reflect.ValueOf(param)
		if paramReflected.Kind() != reflect.Ptr {
			return ErrParamIsNotPointer
		}

		if paramReflected.IsNil() {
			return ErrParamIsNil
		}

		kindOfElement := paramReflected.Elem().Kind()
		if kindOfElement == reflect.Ptr {
			ptr := paramReflected.Elem()
			if ptr.IsNil() {
				ptr.Set(reflect.New(ptr.Type().Elem()))
			}
			param = ptr.Interface()
		} else {
			if kindOfElement == reflect.Slice {
				err = stringToSlice(s, paramReflected)
			} else if kindOfElement == reflect.Array || kindOfElement == reflect.Map {
				err = &ErrUnsupportedKind{kind: kindOfElement}
			} else {
				err = stringToPrimitive(s, paramReflected)
			}
			break
		}
	}

	return
}

var ErrParamIsNotPointer error = errors.New("param isn't a pointer")
var ErrParamIsNil error = errors.New("param is nil")

type ErrUnsupportedKind struct {
	kind reflect.Kind
}

func (err *ErrUnsupportedKind) Error() string {
	return fmt.Sprintf("unsupported kind: '%s'", err.kind.String())
}

type ErrValueIsNotPointer struct {
	value reflect.Value
}

func (err *ErrValueIsNotPointer) Error() string {
	return fmt.Sprintf("value is not a pointer: '%s'", err.value.Kind().String())
}

type ErrTypeValueOverflow struct {
	value string
}

func (err *ErrTypeValueOverflow) Error() string {
	return fmt.Sprintf("type overflow: '%s'", err.value)
}

type ErrUnsupportedPrimitiveType struct {
	value reflect.Value
}

func (err *ErrUnsupportedPrimitiveType) Error() string {
	return fmt.Sprintf("unsupported primitive type: '%s'", err.value.Kind().String())
}

type RecoverError struct {
	Err   interface{}
	Stack []byte
}

func (e *RecoverError) Error() string { return fmt.Sprintf("fromString panicked: %v", e.Err) }
