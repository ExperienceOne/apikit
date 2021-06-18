package middleware

import (
	"errors"
	"fmt"
	"reflect"
)

// ClearFieldByType clears all fields and nested fields in an object
// that have a specified type.
func ClearFieldByType(obj interface{}, t reflect.Type) error {

	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return &ErrValueIsNotPointer{value: value}
	}

	clearValueFieldByType(value, t)

	return nil
}

func clearValueFieldByType(value reflect.Value, t reflect.Type) {

	switch value.Kind() {
	case reflect.Ptr:
		clearValueFieldByType(value.Elem(), t)

	case reflect.Interface:
		if value.Elem().IsValid() && value.Elem().Type() == t {
			value.Set(reflect.Zero(t))
		} else {
			clearValueFieldByType(value.Elem(), t)
		}

	case reflect.Struct:
		if value.Type() == t && value.IsValid() && value.CanSet() {
			value.Set(reflect.Zero(t))
		} else {
			for i := 0; i < value.NumField(); i++ {
				field := value.Field(i)
				if field.Type() == t && field.IsValid() && field.CanSet() {
					field.Set(reflect.Zero(t))
				} else {
					clearValueFieldByType(field, t)
				}
			}
		}

	case reflect.Slice:
		sliceCopy := reflect.New(value.Type()).Elem()
		sliceCopy.Set(reflect.MakeSlice(value.Type(), 0, 0))
		for i := 0; i < value.Len(); i++ {
			sliceValue := value.Index(i)
			if sliceValue.Type() == t || (sliceValue.Kind() == reflect.Interface && sliceValue.Elem().Type() == t) {
				continue
			}

			var sliceValueCopy reflect.Value
			if sliceValue.Kind() == reflect.Interface || sliceValue.Kind() == reflect.Ptr {
				sliceValueCopy = reflect.New(sliceValue.Elem().Type())
				sliceValueCopy.Elem().Set(sliceValue.Elem())
			} else {
				sliceValueCopy = reflect.New(sliceValue.Type())
				sliceValueCopy.Elem().Set(reflect.Indirect(sliceValue))
			}

			clearValueFieldByType(sliceValueCopy, t)

			if sliceValue.Type().Kind() == reflect.Ptr {
				sliceCopy.Set(reflect.Append(sliceCopy, sliceValueCopy))
			} else {
				sliceCopy.Set(reflect.Append(sliceCopy, sliceValueCopy.Elem()))
			}
		}
		value.Set(sliceCopy)

	case reflect.Map:
		for _, key := range value.MapKeys() {
			mapValue := value.MapIndex(key)
			if mapValue.Type() == t {
				mapValue.Set(reflect.Zero(t))
			} else {
				var copyValue reflect.Value
				if mapValue.Kind() == reflect.Interface || mapValue.Kind() == reflect.Ptr {
					copyValue = reflect.New(mapValue.Elem().Type())
					copyValue.Elem().Set(mapValue.Elem())
				} else {
					copyValue = reflect.New(mapValue.Type())
					copyValue.Elem().Set(mapValue.Elem())
				}

				clearValueFieldByType(copyValue, t)
				mapValue = copyValue.Elem()
			}
			value.SetMapIndex(key, mapValue)
		}

	default:
		if value.IsValid() && value.CanSet() && value.Type() == t {
			value.Set(reflect.Zero(t))
		}
	}
}

// ClearFieldByName clears all fields and nested fields in an object
// that have a specified name.
func ClearFieldByName(obj interface{}, fieldName string) error {

	if fieldName == "" {
		return ErrFieldNameIsRequired
	}

	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return &ErrValueIsNotPointer{value: value}
	}

	clearValueFieldByName(value, fieldName)

	return nil
}

func clearValueFieldByName(value reflect.Value, fieldName string) {

	switch value.Kind() {
	case reflect.Ptr:
		clearValueFieldByName(reflect.Indirect(value), fieldName)

	case reflect.Interface:
		clearValueFieldByName(value.Elem(), fieldName)

	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			fieldType := value.Type().Field(i)

			if fieldType.Name == fieldName && field.IsValid() && field.CanSet() {
				field.Set(reflect.Zero(fieldType.Type))
			} else {
				clearValueFieldByName(field, fieldName)
			}
		}

	case reflect.Slice:
		sliceCopy := reflect.New(value.Type()).Elem()
		sliceCopy.Set(reflect.MakeSlice(value.Type(), 0, 0))
		for i := 0; i < value.Len(); i++ {
			sliceValue := value.Index(i)

			var sliceValueCopy reflect.Value
			if sliceValue.Kind() == reflect.Interface || sliceValue.Kind() == reflect.Ptr {
				sliceValueCopy = reflect.New(sliceValue.Elem().Type())
				sliceValueCopy.Elem().Set(sliceValue.Elem())
			} else {
				sliceValueCopy = reflect.New(sliceValue.Type())
				sliceValueCopy.Elem().Set(reflect.Indirect(sliceValue))
			}

			clearValueFieldByName(sliceValueCopy, fieldName)

			if sliceValue.Type().Kind() == reflect.Ptr {
				sliceCopy.Set(reflect.Append(sliceCopy, sliceValueCopy))
			} else {
				sliceCopy.Set(reflect.Append(sliceCopy, sliceValueCopy.Elem()))
			}
		}
		value.Set(sliceCopy)

	case reflect.Map:
		mapCopy := reflect.New(value.Type()).Elem()
		mapCopy.Set(reflect.MakeMap(value.Type()))

		for _, key := range value.MapKeys() {
			originalVal := value.MapIndex(key)
			originalValType := originalVal.Type()
			copyVal := reflect.New(originalValType).Elem()

			if key.String() == fieldName {
				if originalVal.Elem().IsValid() {
					copyVal.Set(reflect.Zero(originalVal.Elem().Type()))
				}
			} else if originalValType.Kind() != reflect.Interface || originalVal.Elem().IsValid() {

				if originalVal.Kind() == reflect.Interface {
					originalVal = originalVal.Elem()
					copyVal = reflect.New(originalVal.Type()).Elem()
				}

				deepCopy(originalVal, copyVal)
				clearValueFieldByName(copyVal, fieldName)
			}

			mapCopy.SetMapIndex(key, copyVal)
		}

		value.Set(mapCopy)
	}
}

func deepCopy(original, copy reflect.Value) {

	switch original.Kind() {
	case reflect.Ptr:
		originalVal := original.Elem()
		if !originalVal.IsValid() {
			return
		}
		copy.Set(reflect.New(originalVal.Type()))
		deepCopy(originalVal, copy.Elem())

	case reflect.Interface:
		if !original.IsNil() {
			originalVal := original.Elem()
			copyVal := reflect.New(originalVal.Type()).Elem()
			deepCopy(originalVal, copyVal)
			copy.Set(copyVal)
		}

	case reflect.Struct:
		for i := 0; i < original.NumField(); i++ {
			deepCopy(original.Field(i), copy.Field(i))
		}

	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			deepCopy(original.Index(i), copy.Index(i))
		}

	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalVal := original.MapIndex(key)
			copyVal := reflect.New(originalVal.Type()).Elem()
			deepCopy(originalVal, copyVal)
			copy.SetMapIndex(key, copyVal)
		}

	default:
		copy.Set(original)
	}
}

var ErrFieldNameIsRequired = errors.New("field name is required")

type ErrValueIsNotPointer struct {
	value reflect.Value
}

func (err *ErrValueIsNotPointer) Error() string {
	return fmt.Sprintf("value is not a pointer: '%s'", err.value.Kind().String())
}
