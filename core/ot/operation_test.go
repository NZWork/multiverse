package ot

import (
	"fmt"
	"testing"
)

func TestShift(t *testing.T) {
	ops := []Operation{
		Insert(2),
		Retain(1),
	}

	fmt.Println(Shift(ops))
	for i, o := range ops {
		fmt.Printf("%d %v\n", i, o)
	}
}
