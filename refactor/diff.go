package refactor

import "github.com/pmezard/go-difflib/difflib"

func diffText(beforeText, afterText string, context int) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(beforeText),
		B:        difflib.SplitLines(afterText),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  context,
	}
	return difflib.GetUnifiedDiffString(diff)
}
