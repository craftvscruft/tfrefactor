package refactor

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pmezard/go-difflib/difflib"
)

func Rename(fromAddressString, toAddressString, configPath string) error {
	configPattern := filepath.Join(configPath, "*.tf")
	_, _ = fmt.Println(configPattern)
	filenames, err := filepath.Glob(configPattern)
	if err != nil {
		return err
	}
	fromAddress := ParseAddress(fromAddressString)
	toAddress := ParseAddress(toAddressString)
	if len(fromAddress.RefNameArray()) != len(toAddress.RefNameArray()) {
		return fmt.Errorf("Addresses are different lengths: '%v' and '%v'", fromAddress.RefName(), toAddress.RefName())
	}
	for _, filename := range filenames {
		parsedFile, err := ParseHclFile(filename)
		beforeText := string(parsedFile.Bytes())
		if err != nil {
			return err
		}
		err = RenameInFile(filename, parsedFile, fromAddress, toAddress)
		if err != nil {
			return err
		}
		afterText := string(parsedFile.Bytes())
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(beforeText),
			B:        difflib.SplitLines(afterText),
			FromFile: "Original",
			ToFile:   "Current",
			Context:  3,
		}
		diffText, _ := difflib.GetUnifiedDiffString(diff)
		if len(diffText) > 0 {
			fmt.Printf("Diff for %v\n%v\n", filename, diffText)
		}

	}
	return nil
}

func RenameInFile(filename string, file *hclwrite.File, fromAddress, toAddress *Address) error {
	matchingBlocks := findBlocks(file.Body(), fromAddress)
	for _, block := range matchingBlocks {
		_, _ = fmt.Printf("Renaming %v %v in %v\n", block.Type(), block.Labels(), filename)
		block.SetType(string(toAddress.BlockType()))
		block.SetLabels(toAddress.labels)
	}
	RenameVariablePrefixInBody("", file.Body(), fromAddress, toAddress)
	return nil
}

func RenameVariablePrefixInBody(blockType string, body *hclwrite.Body, fromAddress, toAddress *Address) {
	for name, attr := range body.Attributes() {
		if !(blockType == "moved" && name == "from") {
			attr.Expr().RenameVariablePrefix(fromAddress.RefNameArray(), toAddress.RefNameArray())
		}
	}
	for _, blk := range body.Blocks() {
		RenameVariablePrefixInBody(blk.Type(), blk.Body(), fromAddress, toAddress)
	}
}

func findBlocks(b *hclwrite.Body, address *Address) []*hclwrite.Block {
	var matched []*hclwrite.Block
	for _, block := range b.Blocks() {
		if string(address.BlockType()) == block.Type() {
			if matchLabels(address.labels, block.Labels()) {
				matched = append(matched, block)
			}
		}
	}

	return matched
}

// matchLabels returns true only if the matched and false otherwise.
func matchLabels(lhs []string, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i := range lhs {
		if !(lhs[i] == rhs[i]) {
			return false
		}
	}

	return true
}
