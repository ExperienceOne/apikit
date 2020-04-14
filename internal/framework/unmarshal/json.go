package unmarshal

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	NullError = errors.New("unexpected null value")
	TypeError = errors.New("unexpected type")
)

func JSON(r io.Reader, v interface{}, required bool) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%+v", r))
		}
	}()

	typ := reflect.TypeOf(v)

	if typ.Kind() != reflect.Ptr {
		return errors.New("please pass pointer to target")
	}

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	decoder := json.NewDecoder(r)

	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map || typ.Kind() == reflect.Struct {

		handleNull := func() error {
			if required {
				return NullError
			} else {
				reflect.ValueOf(v).Elem().Set(reflect.Zero(reflect.TypeOf(v).Elem()))
				return nil
			}
		}

		var err error
		var value reflect.Value

		if typ.Kind() == reflect.Slice {

			var abstractSlice []interface{}
			if err := decoder.Decode(&abstractSlice); err != nil {
				return err
			}

			if abstractSlice == nil {
				return handleNull()
			}

			value, err = slice2concrete(typ, abstractSlice)

		} else {

			var abstractMap map[string]interface{}
			if err := decoder.Decode(&abstractMap); err != nil {
				return err
			}

			if abstractMap == nil {
				return handleNull()
			}

			if typ.Kind() == reflect.Map {
				value, err = map2concrete(typ, abstractMap)
			} else {
				value, err = map2object(typ, abstractMap)
			}
		}

		if err != nil {
			return err
		}

		setValue(value, reflect.ValueOf(v).Elem())

	} else {

		if !required {
			return decoder.Decode(v)
		}

		if err := decoder.Decode(&v); err != nil {
			return err
		}

		if v == nil {
			return NullError
		}
	}

	return nil
}

func setValue(value, dest reflect.Value) {

	for dest.Kind() == reflect.Ptr {
		newDest := reflect.New(dest.Type().Elem())
		dest.Set(newDest)
		dest = newDest.Elem()
	}
	dest.Set(value)
}

func setMapValue(key, value, m reflect.Value) {

	mapValue := reflect.New(m.Type().Elem()).Elem()
	setValue(value, mapValue)
	m.SetMapIndex(key, mapValue)
}

func slice2concrete(typ reflect.Type, s []interface{}) (reflect.Value, error) {

	concretSlice := reflect.MakeSlice(typ, len(s), cap(s))

	for i, value := range s {
		concretValue, err := convert(value, typ.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		setValue(concretValue, concretSlice.Index(i))
	}

	return concretSlice, nil
}

func map2concrete(typ reflect.Type, m map[string]interface{}) (reflect.Value, error) {

	concreteMap := reflect.MakeMap(typ)

	for key, value := range m {

		concreteKey, err := convert(key, typ.Key())
		if err != nil {
			return reflect.Value{}, err
		}

		concreteValue, err := convert(value, typ.Elem())
		if err != nil {
			return reflect.Value{}, err
		}

		setMapValue(concreteKey, concreteValue, concreteMap)
	}

	return concreteMap, nil
}

func map2object(typ reflect.Type, m map[string]interface{}) (reflect.Value, error) {

	object := reflect.New(typ).Elem()

	for i := 0; i < typ.NumField(); i++ {

		field := typ.Field(i)
		if field.Anonymous || !object.Field(i).CanSet() {
			continue
		}

		required := false
		key := field.Name

		jsonTags := strings.Split(field.Tag.Get("json"), ",")
		if len(jsonTags) > 0 {

			if jsonTags[0] == "-" {
				continue
			} else if jsonTags[0] != "" {
				key = jsonTags[0]
			}

			for i := 1; i < len(jsonTags); i++ {
				if jsonTags[i] == "required" {
					required = true
					break
				}
			}
		}

		if value, exists := m[key]; exists && value != nil {

			concreteValue, err := convert(value, field.Type)
			if err != nil {
				return reflect.Value{}, err
			}
			setValue(concreteValue, object.Field(i))

		} else if required {

			return reflect.Value{}, NullError
		}
	}

	return object, nil
}

func convert(data interface{}, typ reflect.Type) (reflect.Value, error) {

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map || typ.Kind() == reflect.Struct {

		var err error
		var value reflect.Value

		if typ.Kind() == reflect.Slice {

			abstractSlice, isSlice := data.([]interface{})
			if !isSlice {
				return reflect.Value{}, TypeError
			}

			value, err = slice2concrete(typ, abstractSlice)

		} else if typ.Kind() == reflect.Map || typ.Kind() == reflect.Struct {

			abstractMap, isMap := data.(map[string]interface{})

			if !isMap {
				return reflect.Value{}, TypeError
			}

			if typ.Kind() == reflect.Map {
				value, err = map2concrete(typ, abstractMap)
			} else {
				value, err = map2object(typ, abstractMap)
			}

		}

		if err != nil {
			return reflect.Value{}, err
		}

		return value, nil

	} else {

		return reflect.ValueOf(data).Convert(typ), nil
	}
}
