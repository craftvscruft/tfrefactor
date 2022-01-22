package refactor

import (
	"fmt"
	"path/filepath"

	"github.com/raymyers/hcl/v2"
	"github.com/raymyers/hcl/v2/hclwrite"
)

func Rename(fromAddressString, toAddressString, configPath string) (*UpdatePlan, error) {
	configPattern := filepath.Join(configPath, "*.tf")
	_, _ = fmt.Println(configPattern)
	filenames, err := filepath.Glob(configPattern)
	if err != nil {
		return nil, err
	}
	fromAddress := ParseAddress(fromAddressString)
	toAddress := ParseAddress(toAddressString)
	if len(fromAddress.RefNameArray()) != len(toAddress.RefNameArray()) {
		return nil, fmt.Errorf("Addresses are different lengths: '%v' and '%v'", fromAddress.RefName(), toAddress.RefName())
	}
	plan := newUpdatePlan()
	for _, filename := range filenames {
		parsedFile, err := ParseHclFile(filename)
		beforeText := string(parsedFile.Bytes())
		if err != nil {
			return nil, err
		}
		err = RenameInFile(filename, parsedFile, fromAddress, toAddress)
		if err != nil {
			return nil, err
		}
		afterText := string(parsedFile.Bytes())
		diffText, err := diffText(beforeText, afterText, 3)
		if len(diffText) > 0 {
			fmt.Printf("Diff for %v\n%v\n", filename, diffText)
			plan.addFileUpdate(&FileUpdate{filename, beforeText, afterText})
		}

	}
	return &plan, nil
}

func createTraversal(labels []string) (traversal hcl.Traversal) {
	traversal = hcl.Traversal{
		hcl.TraverseRoot{
			Name: labels[0],
		},
	}
	for _, label := range labels[1:] {
		traversal = append(traversal, hcl.TraverseAttr{
			Name: label,
		})
	}
	return
}

func RenameInFile(filename string, file *hclwrite.File, fromAddress, toAddress *Address) error {
	if fromAddress.elementType == TypeLocal {
		if err := RenameLocalInFile(filename, file, fromAddress, toAddress); err != nil {
			return err
		}
	} else {
		matchingBlocks := findBlocks(file.Body(), fromAddress)
		for _, block := range matchingBlocks {
			_, _ = fmt.Printf("Renaming %v %v in %v\n", block.Type(), block.Labels(), filename)
			block.SetType(string(toAddress.BlockType()))
			block.SetLabels(toAddress.labels)
			if fromAddress.elementType == TypeResource && toAddress.elementType == TypeResource {
				AddMovedBlock(file, fromAddress, toAddress)
			}
		}
	}

	RenameVariablePrefixInBody("", file.Body(), fromAddress, toAddress)
	return nil
}

func AddMovedBlock(file *hclwrite.File, fromAddress, toAddress *Address) {
	file.Body().AppendNewline()
	movedBlock := file.Body().AppendNewBlock("moved", []string{})

	movedBlock.Body().SetAttributeTraversal("from", createTraversal(fromAddress.labels))
	movedBlock.Body().SetAttributeTraversal("to", createTraversal(toAddress.labels))
}

func RenameLocalInFile(filename string, file *hclwrite.File, fromAddress, toAddress *Address) error {
	fromName := fromAddress.labels[0]
	toName := toAddress.labels[0]
	for _, block := range file.Body().Blocks() {
		if "locals" == block.Type() {
			attr := block.Body().GetAttribute(fromName)
			if attr != nil {
				attr.SetName(toName)
			}
		}
	}
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
		if string(address.BlockType()) == block.Type() && matchLabels(address.labels, block.Labels()) {
			matched = append(matched, block)
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
