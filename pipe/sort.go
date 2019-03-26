package pipe

import (
	"reflect"
	"sort"
)

func SortBy(data interface{}, swap interface{}) interface{} {
	assertFn(swap)

	_fn := reflect.ValueOf(swap)
	_t := _fn.Type()
	assert(_t.NumIn() != 2, "the func input num is more than 2(%d)", _t.NumIn())
	assert(_t.Out(0).Kind() != reflect.Bool, "the func output type is not bool(%s)", _t.Out(0).Kind().String())

	_d := reflect.ValueOf(data)
	var _ps []reflect.Value
	for i := 0; i < _d.Len(); i++ {
		if !_d.Index(i).IsValid() {
			_ps = append(_ps, reflect.Zero(_t.In(0)))
		} else {
			_ps = append(_ps, _d.Index(i))
		}
	}

	_st := reflectValueSlice{data: _ps, swap: _fn}.Sort()
	_rst := reflect.MakeSlice(_d.Type(), 0, _d.Len())
	_rst = reflect.Append(_rst, _st...)

	return _rst.Interface()
}

type reflectValueSlice struct {
	data []reflect.Value
	swap reflect.Value
}

func (c reflectValueSlice) Len() int {
	return len(c.data)
}
func (c reflectValueSlice) Swap(i, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
}

func (c reflectValueSlice) Sort() []reflect.Value {
	sort.Sort(c)
	return c.data
}

func (c reflectValueSlice) Less(i, j int) bool {
	return c.swap.Call([]reflect.Value{c.data[i], c.data[j]})[0].Bool()
}
