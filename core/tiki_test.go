package core

import "testing"

const (
	text1 = "# Intro Hello tiki! 我是谁，你叫什么"
	text2 = "# Intro Hey Tiki! 我也不知道"
)

func TestDiff(t *testing.T) {
	// dmp := diffmatchpatch.New()
	//
	// diffs := dmp.DiffMain(text1, text2, true) // []Diff
	// fmt.Printf("text1: %s\n", text1)
	// fmt.Printf("text2: %s\n", text2)
	// fmt.Println(dmp.DiffPrettyText(diffs))
	// patches := dmp.PatchMake(text1, diffs)
	// fmt.Println("------- result -------")
	// final, _ := dmp.PatchApply(patches, text1)
	// fmt.Println(final)
	// // Patch Apply
}
