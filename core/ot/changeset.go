package ot

import (
	"bytes"
	"errors"
	"log"
)

// Changeset struct
type Changeset struct {
	OP           []Operation `json:"op"`
	Adden        string      `json:"adden"`
	Removen      string      `json:"removen"`
	InputLength  uint        `json:"inputLength"`
	OutputLength uint        `json:"outputLength"`
}

// Apply changeset to text
func (c *Changeset) Apply(content string) (newContent string, err error) {
	var (
		addenPointer, removenPointer, pointer, temp uint
		buf                                         bytes.Buffer
		op                                          Operation
		tempString, tempString1                     string
	)

	for _, op = range c.OP {
		switch op.Type {
		case OPRetain:
			temp = op.Length
			// bce
			tempString, err = UTF8SubString(content, pointer, temp)
			if err != nil {
				return "", err
			}
			buf.WriteString(tempString)
			pointer += temp
		case OPInsert:
			temp = op.Length
			tempString, err = UTF8SubString(c.Adden, addenPointer, temp)
			if err != nil {
				return "", err
			}
			buf.WriteString(tempString)
			addenPointer += temp
		case OPDelete:
			temp = op.Length
			tempString, err = UTF8SubString(content, pointer, temp)
			if err != nil {
				return "", err
			}
			tempString1, err = UTF8SubString(c.Removen, removenPointer, temp)
			if err != nil {
				return "", err
			}

			if tempString == tempString1 {
				pointer += temp
				removenPointer += temp
			} else {
				return "", errors.New("conflict")
			}
		}
	}
	if pointer < uint(UTF8RealLength(content)) {
		tempString, err = UTF8SubString(content, pointer, 0)
		if err != nil {
			return "", err
		}
		buf.WriteString(tempString)
	}

	newContent = buf.String()
	c.OutputLength = uint(UTF8RealLength(newContent))

	return
}

// Invert all the operations in changeset and return it self
func (c *Changeset) Invert() *Changeset {
	for i := 0; i < len(c.OP); i++ {
		c.OP[i] = c.OP[i].Invert()
	}

	temp := c.Adden
	c.Adden = c.Removen
	c.Removen = temp

	return c
}

// IntentionPreservation transform an editing operation into a new form according to the effects of previously executed concurrent operations
// so that the transformed operation can achieve the correct effect
func (c *Changeset) IntentionPreservation(pre *Changeset) (shouldForceSync bool) {
	ops1, ops2 := Clone(pre.OP), Clone(c.OP)
	var op1, op2 Operation
	var tempOp1, tempOp2 Operation
	var newStack []Operation

	debug := false

	for {
		if len(ops1) == 0 {
			if len(ops2) != 0 {
				newStack = append(newStack, ops2...)
				break
			} else {
				break
			}
		}

		if len(ops2) == 0 {
			break
		}
		if debug {
			log.Printf("\n%v\n%v\n", OperationesToString(ops1), OperationesToString(ops2))
		}

		ops1, op1 = Shift(ops1)
		ops2, op2 = Shift(ops2)

		if debug {
			log.Printf("origin op1: %v, op2: %v\n", op1.Debug(), op2.Debug())
		}

		if op1.Length > op2.Length {
			tempOp1, tempOp2 = op1.Derive(op2.Length)
			ops1 = Unshift(ops1, tempOp2)
			op1 = tempOp1
		} else if op2.Length > op1.Length {
			tempOp1, tempOp2 = op2.Derive(op1.Length)
			ops2 = Unshift(ops2, tempOp2)
			op2 = tempOp1
		}
		if debug {
			log.Printf("op1: %v, op2: %v\n", op1.Debug(), op2.Debug())
		}

		if op1.Type == OPRetain {
			if op2.Type == OPRetain {
				newStack = append(newStack, op2)
			} else if op2.Type == OPInsert {
				newStack = append(newStack, op2)
			} else if op2.Type == OPDelete {
				newStack = append(newStack, op2)
			}
		} else if op1.Type == OPInsert {
			if op2.Type == OPRetain {
				op2.Length += op1.Length
				newStack = append(newStack, op2)
			} else if op2.Type == OPInsert || op2.Type == OPDelete {
				newStack = append(newStack, Retain(op1.Length))
				newStack = append(newStack, op2)
			}
		} else if op1.Type == OPDelete {
			if op2.Type == OPRetain {
				// newStack = append(newStack, op2)
			} else if op2.Type == OPInsert || op2.Type == OPDelete {
				backwardPosition := len(newStack) - 1

				temp, inStack := Get(newStack, backwardPosition)
				if !inStack || backwardPosition < 0 {
					shouldForceSync = true
					return
				}

				temp.Length -= op1.Length
				if temp.Length == 0 {
					newStack = newStack[:backwardPosition]
				}
				newStack = append(newStack, op2)
			}
		}
		if debug {
			log.Printf("new %s\n========================\n", OperationesToString(newStack))
		}
	}

	c.OP = newStack
	return
}

/**
 * Inclusion Transformation (IT) or Forward Transformation
 *
 * transforms the operations of the current changeset against the
 * all operations in another changeset in such a way that the
 * effects of the latter are effectively included.
 * This is basically like a applying the other cs on this one.
 *
 * @param otherCs <Changeset>
 * @param left <boolean> Which op to choose if there's an insert tie (If you use this function in a distributed, synchronous environment, be sure to invert this param on the other site, otherwise it can be omitted safely)
 *
 * @returns <Changeset>
 */

// Merge two changesets (that are based on the same state!) so that the resulting changseset
// has the same effect as both orignal ones applied one after the other
func (c *Changeset) Merge(another *Changeset) {

}

// Debug changeset
func (c *Changeset) Debug() {
	for _, o := range c.OP {
		log.Println(o.Debug())
	}
}
