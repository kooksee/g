package pipe

import (
	"github.com/kooksee/g/internal/assert"
	"reflect"
)

func IsError(p interface{}) bool {
	if p == nil {
		return false
	}

	_, ok := p.(error)
	return ok
}

func Type(p interface{}) reflect.Kind {
	return reflect.TypeOf(p).Kind()
}

func Fn(f interface{}, params ...interface{}) func() interface{} {
	return func() interface{} {
		t := reflect.TypeOf(f)
		assert.Bool(t.Kind() != reflect.Func, "err -> Wrap: please input func")

		var vs []reflect.Value
		for i, p := range params {
			if p == nil {
				vs = append(vs, reflect.New(t.In(i)).Elem())
			} else {
				vs = append(vs, reflect.ValueOf(p))
			}
		}

		out := reflect.ValueOf(f).Call(vs)
		if !out[0].IsValid() {
			return nil
		}

		return out[0]
	}
}

func assertFn(fn interface{}) {
	assert.Bool(fn == nil, "the func is nil")

	_v := reflect.ValueOf(fn)
	assert.Bool(_v.Kind() != reflect.Func, "the params(%s) is not func type", _v.Type())
}
