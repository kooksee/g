package pipe_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kooksee/pipe"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

type t1 struct {
	A string
	b int
}

func TestP(t *testing.T) {
	pipe.Data([]int{1, 2, 3}, []int{1, 2, 3}).P()
	pipe.Data(t1{A: "dd", b: 1}, &t1{A: "sss", b: 2}).P()
}

func TestFilter(t *testing.T) {
	t.Run("test filter", func(t *testing.T) {
		pipe.Data(t1{A: "dd", b: 1}, &t1{A: "sss", b: 2}).Filter(func(i int, v interface{}) bool {
			return !pipe.IsPtr(v)
		}).P()
	})

	t.Run("test filter type", func(t *testing.T) {
		pipe.Data(&t1{A: "dd", b: 1}, &t1{A: "sss", b: 2}).Filter(func(v *t1) bool {
			return v.b > 1
		}).P()
	})
}

func TestMap(t *testing.T) {
	pipe.Data(t1{A: "dd", b: 1}, t1{A: "sss", b: 2}).Map(func(i int, v t1) t1 {
		v.b = 100000
		return v
	}).P("test map")

	pipe.Data(t1{A: "dd", b: 1}, t1{A: "sss", b: 2}).Map(func(i int, v t1) t1 {
		v.b = 100000
		return v
	}).Each(func(i int, a ...t1) {
		fmt.Println(a)
	})

	pipe.Data(t1{A: "dd", b: 1}, t1{A: "sss", b: 2}).Map(func(i int, v t1) t1 {
		v.b = 100000
		return v
	}).Pipe(func(a t1, aa t1) {
		fmt.Println(a)
	})

	pipe.Data(nil, &t1{}).Map(func(i int, v *t1) *t1 {
		if v == nil {
			return nil
		}

		fmt.Println(v.b)

		v.b = 100000

		return v
	}).Map(func(v *t1) *t1 {
		if v == nil {
			return nil
		}

		fmt.Println("map2", v.b)
		v.b = 222000000
		return v
	}).Each(func(v *t1) {
		fmt.Println(v)
	})
}

func TestArray(t *testing.T) {
	var ddd []int
	ddd = append(ddd, 1, 2, 34)
	pipe.DataArray(ddd).Each(func(i, n int) {
		fmt.Println(i, n)
	})

	pipe.DataRange(1, 100, 3).P()
	pipe.DataRange(1, 100, 3).Map(func(a int) int {
		fmt.Println(a)
		return a
	}).Each(func(i, n int) {
		fmt.Println(i, n)
	})
}

func TestReduce(t *testing.T) {

	pipe.Data(t1{A: "dd", b: 1}, t1{A: "sss", b: 2}).Map(func(v t1) t1 {
		v.b = 100000
		return v
	}).Reduce(func(s t1, v t1) t1 {
		return t1{b: s.b + v.b, A: s.A + v.A}
	}).Each(func(a interface{}) {
		fmt.Println(a)
	})

	fmt.Println(pipe.Data(t1{A: "dd", b: 1}, t1{A: "sss", b: 2}, t1{A: "sss", b: 2}).Map(func(i int, v t1) t1 {
		v.b = 100000
		return v
	}).Reduce(func(s t1, v t1) t1 {
		return t1{b: s.b + v.b, A: s.A + v.A}
	}).ToData().Json())
}

func TestEach(t *testing.T) {
	pipe.Data(1, 2, 3, 4).Each(func(a ...interface{}) {
		fmt.Println(a)
	})

	pipe.Data(1, 2, 3, 4).Each(func(i int, a interface{}) {
		fmt.Println(i, a)
	})

	pipe.Data(1, 2, 3, 4).Each(func(a interface{}) {
		fmt.Println(a)
	})
}

func TestPipe(t *testing.T) {
	pipe.Data(1, "dd").Pipe(func(a int, b string) (int, string) {
		fmt.Println("callback success ok", a, b)
		return a, b
	}).Pipe(func(a int, b string) {
		fmt.Println("callback ", a, b)
	})

	pipe.Data(1, 2, 3, 4, nil).Pipe(func(a, b, c, d int, e error) {
		fmt.Println(a, b, c, d, e)
	}).P("test pipe")
}

func TestIsError(t *testing.T) {
	fmt.Println(pipe.IsError(errors.New("")))
	fmt.Println(pipe.IsError(nil))

}

func TestError(t *testing.T) {
	//pipe.Data(1, 2, 3, errors.New("sss")).MustNotError()
	pipe.Data(1, 2, 3, nil).MustNotError()

}

func TestToData(t *testing.T) {
	a := "0 */2 * * * *"
	fmt.Println(pipe.Data(strings.Split(a, "*")[1]).ToData().String())
	fmt.Println(pipe.DataFromArray(strings.Split(a, "*")).ToData().String())
	fmt.Println(pipe.DataFromArray(strings.Split(a, "*")).ToData().Json())
	fmt.Println(pipe.Data(1, 2, 3, "", nil, &a).ToData().Json())
	pipe.Data(1, 2, 3, "", nil, &a).P()
}

func TestIf(t *testing.T) {
	t.Run("懒加载", func(t *testing.T) {
		a := "0 */2 * * * *"
		fmt.Println(pipe.If(true, pipe.Fn(strings.Split, a, "*"), 2))
		fmt.Println(pipe.If(false, pipe.Fn(fmt.Println, "1", 2), 2))
		fmt.Println(pipe.If(true, pipe.Fn(fmt.Println, "1", 2), 2))
	})
}

func TestSetInterface(t *testing.T) {
	_fn := func(in interface{}, a interface{}) {
		fmt.Println(in, a)
		reflect.ValueOf(in).Elem().Set(reflect.ValueOf(a))
	}
	a := 1
	b := 2
	_fn(&a, b)
	fmt.Println(a, b)
}

func TestExpr(t *testing.T) {
	pipe.Data(1, 2, 3, 4, nil).Pipe(func(a, b, c, d int, e error) {
		fmt.Println(a, b, c, d, e)
	}).P("test pipe")

	fmt.Println(pipe.Data(1, 2, 3, 4, nil).FilterExp(`it == nil`).ToData().Json())
	fmt.Println(pipe.Data(&t1{A: "1", b: 2}, &t1{A: "1", b: 3}).FilterExp(`it.A == "1"`).ToData().Json())
	fmt.Println(pipe.Data(&t1{A: "1", b: 2}, &t1{A: "1", b: 3}).MapExp(`it.A == "1"`).ToData().Json())

	_a := pipe.
		Data(nil, &t1{A: "1", b: 2}, &t1{A: "1", b: 3}).
		ToData().
		Interface().([]*t1)
	fmt.Println(_a)
	fmt.Println(_a[1].A)
}

type M struct {
	A  string `json:"a"`
	A1 string `json:"a1"`
	A2 string `json:"a2"`
}

func (t *M) Name() string {
	return "m"
}

func TestSortBy(t *testing.T) {
	if a, ok := pipe.SortBy([]string{"11", "2", "3"}, func(a, b string) bool {
		return strings.Compare(a, b) > 0
	}).([]string); ok {
		fmt.Println(a)
	}

	aa := pipe.SortBy([]string{"11", "2", "31"}, func(a, b string) bool {
		return strings.Compare(a, b) < 0
	}).([]string)

	for _, _a := range aa {
		fmt.Println(_a, len(_a), len(aa))
	}

	fmt.Println("nil test", pipe.
		Data(nil, &t1{A: "1", b: 2}, &t1{A: "1", b: 3}).
		SortBy(func(a, b *t1) bool {
			if a == nil || b == nil {
				return true
			}

			return a.b > b.b
		}).ToData().Interface())

	fmt.Println(pipe.Data(1, 11, 2).SortBy(func(a, b int) bool {
		return a > b
	}).ToData().Json())

	fmt.Println(pipe.Data(1, 11, 2).SortBy(func(a, b int) bool {
		return a > b
	}).ToData().Interface())
}

func TestGroupBy(t *testing.T) {
	pipe.Data(nil, &t1{A: "1", b: 2}, map[string]interface{}{"A": "2"}, t1{A: "1", b: 2}, &t1{A: "2", b: 3}, &t1{A: "2", b: 3}).GroupBy("A").Echo()
}

func TestName123(t *testing.T) {
	a := "1a38753878917df87acfac"
	b, err := hex.DecodeString(a)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(big.NewInt(0).SetBytes(b).Int64())
}
