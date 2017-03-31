package ot

import (
	"fmt"
	"log"
	"testing"
)

func TestApply(t *testing.T) {
	c := Changeset{
		OP: []Operation{
			Retain(1),
			Insert(1),
			Retain(1),
			Insert(1),
			Retain(3),
			Insert(1),
			Delete(1),
			Retain(1),
			Insert(1),
		},
		Adden:        "elo!",
		Removen:      "t",
		InputLength:  7,
		OutputLength: 10,
	}

	content, err := c.Apply("hlo, tt")
	if err != nil {
		t.Fatal(err)
	}

	if content != "hello, ot!" {
		t.Fatalf("given %s\n", content)
	}
}

func TestIntentionPreservation(t *testing.T) {
	fmt.Println("intention preservation")
	c1 := &Changeset{
		OP: []Operation{
			Delete(1),
			Retain(1),
			Insert(1),
			Retain(1),
			Delete(1),
			Retain(5),
			Insert(1),
		},
		Adden:        "26",
		Removen:      "14",
		InputLength:  9,
		OutputLength: 9,
	}

	c2 := &Changeset{
		OP: []Operation{
			Retain(2),
			Delete(1),
			Retain(1),
			Delete(1),
			Retain(4),
			Insert(1),
		},
		Adden:        "0",
		Removen:      "35",
		InputLength:  9,
		OutputLength: 8,
	}

	c2.IntentionPreservation(c1)
	log.Println(OperationesToString(c2.OP))
	// c2 => Retain 1 | Retain 1 | Delete 1 | Delete 1 | Retain 4 | Insert 1

	c3 := &Changeset{
		OP: []Operation{
			Delete(1),
			Retain(1),
			Insert(1),
			Delete(2),
			Retain(1),
		},
		Adden:        "4",
		Removen:      "157",
		InputLength:  5,
		OutputLength: 3,
	}

	c4 := &Changeset{
		OP: []Operation{
			Retain(1),
			Delete(1),
			Retain(2),
			Delete(1),
		},
		Adden:        "",
		Removen:      "39",
		InputLength:  5,
		OutputLength: 3,
	}

	c4.IntentionPreservation(c3)
	log.Println(OperationesToString(c4.OP))
	// c4 => Delete 1 | Retain 1 | Delete 1

	// c5 := &Changeset{
	// 	OP: []Operation{
	// 		Delete(1),
	// 		Retain(1),
	// 		Delete(1),
	// 		Retain(1),
	// 		Delete(1),
	// 	},
	// 	Adden:        "",
	// 	Removen:      "159",
	// 	InputLength:  5,
	// 	OutputLength: 2,
	// }
	//
	// c6 := &Changeset{
	// 	OP: []Operation{
	// 		Retain(1),
	// 		Insert(1),
	// 		Delete(1),
	// 		Retain(1),
	// 		Insert(1),
	// 		Delete(1),
	// 		Retain(1),
	// 		Insert(1),
	// 	},
	// 	Adden:        "sss",
	// 	Removen:      "37",
	// 	InputLength:  5,
	// 	OutputLength: 5,
	// }
	//
	// c6.IntentionPreservation(c5)
	// log.Println(OperationesToString(c6.OP))

}
