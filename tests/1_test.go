package tests

import (
	"github.com/kooksee/g/pp"
	"testing"
)

type ss struct {
	A string
}

func TestName(t *testing.T) {
	pp.Println(ss{A: "ddd"})
}
