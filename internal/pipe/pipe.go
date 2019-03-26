package pipe

import (
	"fmt"
	"github.com/antonmedv/expr"
	"reflect"
)

type _func struct {
	params []reflect.Value
}

func (t *_func) SortByName(name string, desc bool) *_func {

	//_fn := func(a, b reflect.Value) {
	//
	//}

	//for i := 0; i < len(t.params); i++ {
	//	if !t.params[i].IsValid() {
	//		t.params[i] = reflect.Zero(_t.In(0))
	//	}
	//}

	//t.params = reflectValueSlice{data: t.params, swap: _fn}.Sort()

	return nil
}

func (t *_func) SortBy(swap interface{}) *_func {
	assertFn(swap)

	_fn := reflect.ValueOf(swap)
	_t := _fn.Type()
	assert(_t.NumIn() != 2, "the func input num is more than 2(%d)", _t.NumIn())
	assert(_t.Out(0).Kind() != reflect.Bool, "the func output type is not bool(%s)", _t.Out(0).Kind().String())

	for i := 0; i < len(t.params); i++ {
		if !t.params[i].IsValid() {
			t.params[i] = reflect.Zero(_t.In(0))
		}
	}

	t.params = reflectValueSlice{data: t.params, swap: _fn}.Sort()

	return t
}

func (t *_func) Pipe(fn interface{}) *_func {
	assertFn(fn)

	_fn := reflect.ValueOf(fn)
	_t := _fn.Type()

	assert(len(t.params) != _t.NumIn(), "the params num is not match(%d,%d)", len(t.params), _t.NumIn())

	var _res []reflect.Value
	for i, p := range t.params {
		if !p.IsValid() {
			p = reflect.New(_t.In(i)).Elem()
		}
		_res = append(_res, p)
	}
	return &_func{params: _fn.Call(_res)}
}

func (t *_func) P(tags ...string) {
	for _, p := range t.params {
		if p.IsValid() {
			fmt.Println(p.Kind().String(), p.Type().String(), p.Interface())
			continue
		}

		fmt.Println("InValid", true)
	}

	if len(tags) > 0 {
		fmt.Println(tags[0])
	}
	fmt.Print("\n\n")
}

func (t *_func) ToData() *_data {
	return &_data{_values: t.params}
}

func (t *_func) MapExp(code string) *_func {
	_ps := make(map[string]interface{})

	var vs []reflect.Value
	for _, p := range t.params {
		if !p.IsValid() {
			_ps["it"] = nil
		} else {
			_ps["it"] = p.Interface()
		}

		output, err := expr.Eval(code, _ps)
		assert(err != nil, "FilterExp error(%s)", err)

		out, ok := output.(bool)
		assert(!ok, "output type error")

		if out {
			vs = append(vs, p)
		}
	}

	return &_func{params: vs}
}

func (t *_func) Map(fn interface{}) *_func {
	assertFn(fn)

	_fn := reflect.ValueOf(fn)
	_t := _fn.Type()
	assert(_t.NumIn() > 2 || _t.NumIn() == 0, "the func input num is [1,2], now(%d)", _t.NumIn())
	assert(_t.NumOut() != 1, "the func output num is 1, now(%d)", _t.NumOut())
	assert(_t.In(_t.NumIn()-1) != _t.Out(0), "the func input output type is not match(%s,%s)", _t.In(_t.NumIn()-1), _t.Out(0))

	var _res []reflect.Value
	for i, p := range t.params {
		if !p.IsValid() {
			p = reflect.New(_t.In(_t.NumIn() - 1).Elem())
		}

		_r := _fn.Call(If(_t.NumIn() == 1, []reflect.Value{p}, []reflect.Value{reflect.ValueOf(i), p}).([]reflect.Value))
		if !_r[0].IsValid() {
			_r[0] = reflect.New(_t.Out(0).Elem())
		}

		_res = append(_res, _r[0])
	}

	return &_func{params: _res}
}

func (t *_func) Reduce(fn interface{}) *_func {
	assertFn(fn)

	_fn := reflect.ValueOf(fn)
	_t := _fn.Type()
	assert(_t.NumIn() != 2, "the func input num is not equal 2(%d)", _t.NumIn())
	assert(_t.NumOut() != 1, "the func output num is not equal 1(%d)", _t.NumOut())
	assert(_t.In(0) != _t.In(1) || _t.In(1) != _t.Out(0), "the func input and output type is not match(%s,%s,%s)", _t.In(0), _t.In(1), _t.Out(0))

	if len(t.params) == 0 {
		return &_func{}
	}

	_tp := reflect.New(_t.In(0)).Elem()
	if len(t.params) < 2 {
		if !t.params[0].IsValid() {
			t.params[0] = _tp
		}
		return &_func{params: t.params}
	}

	if len(t.params) < 3 {
		if !t.params[0].IsValid() {
			t.params[0] = _tp
		}
		if !t.params[1].IsValid() {
			t.params[1] = _tp
		}
		return &_func{params: _fn.Call([]reflect.Value{t.params[0], t.params[1]})}
	}

	_res := _fn.Call([]reflect.Value{t.params[0], t.params[1]})
	for i := 2; i < len(t.params); i++ {
		if !t.params[i].IsValid() {
			t.params[i] = _tp
		}
		_res = _fn.Call([]reflect.Value{_res[0], t.params[i]})
	}
	return &_func{params: _res}
}

func (t *_func) Any(fn func(v interface{}) bool) bool {
	assertFn(fn)

	for _, p := range t.params {
		if fn(If(!p.IsValid(), nil, Fn(p.Interface))) {
			return true
		}
	}
	return false
}

func (t *_func) Every(fn func(v interface{}) bool) bool {
	assertFn(fn)

	for _, p := range t.params {
		if !fn(If(!p.IsValid(), nil, Fn(p.Interface))) {
			return false
		}
	}
	return true
}

func (t *_func) MustNotError() {
	for _, p := range t.params {
		if !p.IsValid() {
			continue
		}

		if d, ok := p.Interface().(error); ok {
			panic(d.Error())
		}
	}
}

func (t *_func) FilterError() *_func {
	return t.Filter(func(v interface{}) bool {
		return !IsError(v)
	})
}

func (t *_func) FilterExp(filter string) *_func {
	_ps := make(map[string]interface{})

	var vs []reflect.Value
	for _, p := range t.params {
		if !p.IsValid() {
			_ps["it"] = nil
		} else {
			_ps["it"] = p.Interface()
		}

		output, err := expr.Eval(filter, _ps)
		assert(err != nil, "FilterExp error(%s)", err)

		out, ok := output.(bool)
		assert(!ok, "output type error")

		if out {
			vs = append(vs, p)
		}
	}

	return &_func{params: vs}
}

func (t *_func) Filter(fn interface{}) *_func {
	assertFn(fn)

	_fn := reflect.ValueOf(fn)
	_t := _fn.Type()
	assert(_t.NumIn() > 2, "the func input num is more than 2(%d)", _t.NumIn())
	assert(_t.NumOut() != 1, "the func output num is not equal 1(%d)", _t.NumOut())
	assert(_t.Out(0).Kind() != reflect.Bool, "the func output type is not bool(%s)", _t.Out(0).Kind().String())

	var vs []reflect.Value
	for i, p := range t.params {
		if !p.IsValid() {
			p = reflect.New(_t.In(_t.NumIn() - 1).Elem())
		}

		_r := _fn.Call(If(_t.NumIn() == 1, []reflect.Value{p}, []reflect.Value{reflect.ValueOf(i), p}).([]reflect.Value))
		if _r[0].Bool() {
			vs = append(vs, p)
		}
	}

	return &_func{params: vs}
}

func (t *_func) ToSlice() *_func {
	var _ps []reflect.Value
	if len(t.params) == 0 {
		return t
	}

	_ds := t.params[0]
	for i := 0; i < _ds.Len(); i++ {
		_ps = append(_ps, _ds.Index(i))
	}
	t.params = _ps
	return t
}

func (t *_func) Each(fn interface{}) {
	assertFn(fn)

	_fn := reflect.ValueOf(fn)
	_t := _fn.Type()
	assert(_t.NumIn() > 2, "the func input num is more than 2(%d)", _t.NumIn())
	assert(_t.NumIn() == 0, "the func input num is more than 2(%d)", _t.NumIn())
	assert(_t.NumOut() != 0, "the func output num is not equal 0(%d)", _t.NumOut())

	for i, p := range t.params {
		if !p.IsValid() {
			p = reflect.New(_t.In(_t.NumIn() - 1).Elem())
		}
		_fn.Call(If(_t.NumIn() == 1, []reflect.Value{p}, []reflect.Value{reflect.ValueOf(i), p}).([]reflect.Value))
	}
}

func DataRange(s, e, t int) *_func {
	assert(s > e, "")

	var _ps []reflect.Value
	for i := s; i < e; i += t {
		_ps = append(_ps, reflect.ValueOf(i))
	}
	return &_func{
		params: _ps,
	}
}

func DataFromArray(ps interface{}) *_func {
	return DataArray(ps)
}

func DataArray(ps interface{}) *_func {
	_d := reflect.ValueOf(ps)
	var _ps []reflect.Value
	for i := 0; i < _d.Len(); i++ {
		_ps = append(_ps, _d.Index(i))
	}
	return &_func{
		params: _ps,
	}
}

func Data(ps ...interface{}) *_func {
	var vs []reflect.Value
	for _, v := range ps {
		vs = append(vs, reflect.ValueOf(v))
	}
	return &_func{
		params: vs,
	}
}
