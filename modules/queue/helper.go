package queue

import (
	"encoding/json"
	"reflect"
)

func toConfig(exemplar, cfg interface{}) (interface{}, error) {
	if reflect.TypeOf(cfg).AssignableTo(reflect.TypeOf(exemplar)) {
		return cfg, nil
	}

	configBytes, ok := cfg.([]byte)
	if !ok {
		configStr, ok := cfg.(string)
		if !ok {
			return nil, ErrInvalidConfiguration{cfg: cfg}
		}

		configBytes = []byte(configStr)
	}

	newVal := reflect.New(reflect.TypeOf(exemplar))
	if err := json.Unmarshal(configBytes, newVal.Interface()); err != nil {
		return nil, ErrInvalidConfiguration{cfg: cfg}
	}

	return newVal.Elem().Interface(), nil
}

func assignableTo(d interface{}, exemplar interface{}) bool {
	return reflect.TypeOf(d).AssignableTo(reflect.TypeOf(exemplar))
}
