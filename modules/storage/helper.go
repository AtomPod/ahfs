package storage

import (
	"encoding/json"
	"reflect"
)

type Marshaler interface {
	Unmarshal(v interface{}) error
}

func ToConfig(exampler, cfg interface{}) (interface{}, error) {

	if reflect.TypeOf(cfg).AssignableTo(reflect.TypeOf(exampler)) {
		return cfg, nil
	}

	if mashaler, ok := cfg.(Marshaler); ok {
		newVal := reflect.New(reflect.TypeOf(exampler))
		if err := mashaler.Unmarshal(newVal.Interface()); err != nil {
			return nil, ErrInvalidConfiguration{config: cfg, err: err}
		}
		return newVal.Elem().Interface(), nil
	}

	configBytes, ok := cfg.([]byte)
	if !ok {
		var configStr string
		configStr, ok = cfg.(string)
		configBytes = []byte(configStr)
	}

	if !ok {
		return nil, ErrInvalidConfiguration{config: cfg}
	}

	newVal := reflect.New(reflect.TypeOf(exampler))
	if err := json.Unmarshal(configBytes, newVal.Interface()); err != nil {
		return nil, ErrInvalidConfiguration{config: cfg, err: err}
	}

	return newVal.Elem().Interface(), nil

}
