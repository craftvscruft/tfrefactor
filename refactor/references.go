package refactor

import (
	"strings"

	"github.com/raymyers/hcl/v2/hclwrite"
)

func findReferencingExpresssions(b *hclwrite.Body, address *Address) []*hclwrite.Expression {
	var matched []*hclwrite.Expression
	addReferencingExpresssions(b, address, &matched)
	return matched
}

func addReferencingExpresssions(body *hclwrite.Body, address *Address, matched *[]*hclwrite.Expression) {
	addressRef := address.RefNameArray()
	for _, attr := range body.Attributes() {
		expr := attr.Expr()
		for _, varTrav := range expr.Variables() {
			travLabelsToMatch := traversalLabels(varTrav)
			println("====")
			println(strings.Join(travLabelsToMatch, "::"))
			println(strings.Join(addressRef, "::"))
			if len(travLabelsToMatch) > len(addressRef) {
				travLabelsToMatch = travLabelsToMatch[0:len(addressRef)]
			}
			if matchLabels(addressRef, travLabelsToMatch) {
				*matched = append(*matched, expr)
			}
		}
	}
	for _, block := range body.Blocks() {
		addReferencingExpresssions(block.Body(), address, matched)
	}
}

func traversalLabels(trav *hclwrite.Traversal) []string {
	// hclwrite currently has almost no API for Traversal, working around.
	buf := strings.Builder{}
	trav.BuildTokens(nil).WriteTo(&buf)
	// This won't work in all cases, like arrays.
	return strings.Split(strings.TrimSpace(buf.String()), ".")
}
