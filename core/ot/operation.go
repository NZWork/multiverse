package ot

const (
	OPRetain = iota
	OPInsert
	OPDelete
)

// Operation op
type Operation struct {
	Type      uint        `json:"type"`
	Operation interface{} `json:"op"`
}

func (o Operation) Equal(other Operation) bool {
	if o.Type == other.Type {
		switch o.Type {
		case OPRetain:
			return o.Operation.(float64) == o.Operation.(float64)
		case OPInsert:
		case OPDelete:
			return o.Operation.(string) == o.Operation.(string)
		}
	}
	return false
}

func Retain(length uint64) Operation {
	return Operation{
		Type:      OPRetain,
		Operation: length,
	}
}

func Insert(content string) Operation {
	return Operation{
		Type:      OPInsert,
		Operation: content,
	}
}

func Delete(content string) Operation {
	return Operation{
		Type:      OPDelete,
		Operation: content,
	}
}
