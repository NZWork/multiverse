package ot

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
