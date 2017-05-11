package ot

import "fmt"

const (
	//OPRetain 0
	OPRetain = iota
	// OPInsert 1
	OPInsert
	// OPDelete 2
	OPDelete
)

type Operation struct {
	Type    uint   `json:"type"`
	Length  uint   `json:"length"`
	Input   uint   `json:"input"`
	Output  uint   `json:"output"`
	Content string `json:"content"`
}

// NewDelete return delete instance
func Delete(length uint) Operation {
	return Operation{
		Type:   OPDelete,
		Length: length,
		Input:  length,
		Output: 0,
	}
}

// NewInsert return insert instance
func Insert(lenth uint) Operation {
	return Operation{
		Type:   OPInsert,
		Length: lenth,
		Input:  0,
		Output: lenth,
	}
}

// NewRetain return retain instance
func Retain(length uint) Operation {
	return Operation{
		Type:   OPRetain,
		Length: length,
		Input:  length,
		Output: length,
	}
}

// Get op from operations by index
func GetClone(ops []Operation, index int) (o Operation, ok bool) {
	if index < len(ops) {
		o = ops[index]
		ok = true
	}
	return
}

func Get(ops []Operation, index int) (o *Operation, ok bool) {
	if index < len(ops) {
		o = &ops[index]
		ok = true
	}
	return
}

func Shift(ops []Operation) (os []Operation, op Operation) {
	switch len(ops) {
	case 0:
		return
	case 1:
		op = ops[0]
		return
	default:
		return ops[1:], ops[0]
	}
}

func Unshift(ops []Operation, op Operation) []Operation {
	temp := []Operation{}
	temp = append(temp, op)
	temp = append(temp, ops...)
	return temp
}

func Clone(ops []Operation) (cloned []Operation) {
	for _, op := range ops {
		cloned = append(cloned, op.Clone())
	}
	return
}

func (o Operation) Clone() (s Operation) {
	switch o.Type {
	case OPRetain:
		return Retain(o.Length)
	case OPInsert:
		return Insert(o.Length)
	case OPDelete:
		return Delete(o.Length)
	}
	return
}

// Equal with another operation
func (o Operation) Equal(another Operation) bool {
	if o.Type == another.Type {
		return o.Input == another.Input && o.Length == another.Length && o.Output == another.Output
	}
	return false
}

func (o Operation) Derive(length uint) (left, right Operation) {
	switch o.Type {
	case OPRetain:
		return Retain(length), Retain(o.Length - length)
	case OPInsert:
		return Insert(length), Insert(o.Length - length)
	case OPDelete:
		return Delete(length), Delete(o.Length - length)
	}
	return
}

func (o Operation) Invert() Operation {
	switch o.Type {
	case OPRetain:
		return o
	case OPInsert:
		return Delete(o.Length)
	case OPDelete:
		return Insert(o.Length)
	}
	return o
}

func (o Operation) Debug() string {
	t := ""
	switch o.Type {
	case OPRetain:
		t = "Retain"
	case OPInsert:
		t = "Insert"
	case OPDelete:
		t = "Delete"
	}
	return fmt.Sprintf("%s %d", t, o.Length)
}
