package pipe

import (
	"fmt"
	"reflect"
	"runtime"
)

type _groupBy struct {
	_data map[interface{}][]reflect.Value
}

func (t *_groupBy) Echo() {
	for i, _p := range t._data {
		fmt.Println(i, _p)
	}
}

func (t *_groupBy) GC() {
	t._data = nil
	runtime.GC()
}

//func (t *_groupBy) Count(name ...string) map[interface{}][]map[string]int {
//	_dt := make(map[interface{}][]map[string]int)
//	for i, _p := range t._data {
//		_dt[i] = len(_p)
//	}
//	return _dt
//}

func (t *_groupBy) Min(name ...string) map[interface{}]int {
	_dt := make(map[interface{}]int)
	for i, _p := range t._data {
		_dt[i] = len(_p)
	}
	return _dt
}

func (t *_groupBy) Max(name ...string) map[interface{}]int {
	_dt := make(map[interface{}]int)
	for i, _p := range t._data {
		_dt[i] = len(_p)
	}
	return _dt
}

func (t *_func) GroupBy(name string) *_groupBy {
	_data := make(map[interface{}][]reflect.Value)
	for _, _p := range t.params {
		if !_p.IsValid() {
			continue
		}

		var _v reflect.Value
		switch _p.Kind() {
		case reflect.Ptr:
			_v = _p.Elem().FieldByName(name)
		case reflect.Struct:
			_v = _p.FieldByName(name)
		case reflect.Map:
			_v = _p.MapIndex(reflect.ValueOf(name))
		default:
			assert(true, "type(%s) error", _p.Kind().String())
		}

		if !_v.IsValid() {
			continue
		}

		_data[_v.Interface()] = append(_data[_v.Interface()], _p)
	}

	return &_groupBy{_data: _data}
}
