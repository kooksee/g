package pipe

import (
	"encoding/json"
	"reflect"
)

type _data struct {
	_values []reflect.Value
}

func (t *_data) String() string {
	if len(t._values) < 1 || !t._values[0].IsValid() || t._values[0].Kind() != reflect.String {
		return ""
	}
	return t._values[0].String()
}

func (t *_data) Raw() []reflect.Value {
	return t._values
}

func (t *_data) Interface() interface{} {
	if len(t._values) < 1 {
		return nil
	}

	var _t reflect.Type
	for _, _v := range t._values {
		if _v.IsValid() {
			_t = _v.Type()
			break
		}
	}

	if _t == nil {
		return nil
	}

	for i := 0; i < len(t._values); i++ {
		if !t._values[i].IsValid() {
			t._values[i] = reflect.Zero(_t)
		}
	}

	_rst := reflect.MakeSlice(reflect.SliceOf(_t), 0, len(t._values))
	_rst = reflect.Append(_rst, t._values...)
	return _rst.Interface()
}

func (t *_data) Json() string {
	var _res []interface{}
	for _, _p := range t._values {
		if !_p.IsValid() {
			_res = append(_res, nil)
		} else {
			_res = append(_res, _p.Interface())
		}
	}

	dt, err := json.Marshal(_res)
	assert(err != nil, "data json error(%s)", err)

	return string(dt)
}
