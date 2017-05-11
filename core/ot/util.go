package ot

import "errors"

// ChangesetFromDiff return a changeset transforms from a diff
func ChangesetFromDiff() (c Changeset) {
	return
}

func OperationesToString(ops []Operation) (s string) {
	for _, op := range ops {
		s += op.Debug()
		s += " | "
	}
	return
}

func UTF8SubString(base string, start, length uint) (string, error) {
	rs := []rune(base)
	if length == 0 && start > 0 {
		return string(rs[start:]), nil
	}
	if length > uint(len(rs))-start || start < 0 {
		return "", errors.New("slice bounds out of range")
	}
	return string(rs[start : start+length]), nil
}

func UTF8RealLength(base string) int {
	return len([]rune(base))
}
