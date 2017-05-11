package ot

import "testing"

func TestUTF8(t *testing.T) {
	temp := "测试utf-8"
	t.Log(UTF8SubString(temp, 0, 3))
}
