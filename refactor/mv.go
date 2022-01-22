package refactor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/raymyers/hcl/v2/hclwrite"
)

func Mv(fromAddressString, toFile, configPath string) (*UpdatePlan, error) {
	filenames, err := filepath.Glob(configPath + "/*.tf")
	if err != nil {
		return nil, err
	}
	plan := newUpdatePlan()
	var parsedOutFile *hclwrite.File

	// Assume toFile is relative to config path. Will this always be true?
	toFilePath := filepath.Join(configPath, toFile)
	if _, err := os.Stat(toFilePath); errors.Is(err, os.ErrNotExist) {
		parsedOutFile, err = ParseHclBytes([]byte{}, toFilePath)
		if err != nil {
			return nil, err
		}
	} else {
		parsedOutFile, err = ParseHclFile(toFilePath)
		if err != nil {
			return nil, err
		}
	}

	beforeOutText := string(parsedOutFile.Bytes())
	for _, filename := range filenames {
		fromPath, _ := filepath.Abs(filename)
		if fromPath != "" && fromPath != toFilePath {
			parsedInFile, err := ParseHclFile(filename)
			if err != nil {
				return nil, err
			}
			beforeText := string(parsedInFile.Bytes())
			if err != nil {
				return nil, err
			}
			fromAddress := ParseAddress(fromAddressString)
			moveAddrToFile(fromAddress, parsedInFile, parsedOutFile)
			afterText := string(parsedInFile.Bytes())
			if err != nil {
				return nil, err
			}
			diffText, err := diffText(beforeText, afterText, 3)
			if len(diffText) > 0 {
				fmt.Printf("Diff for %v\n%v\n", filename, diffText)
				plan.addFileUpdate(&FileUpdate{filename, beforeText, afterText})
			}
		}
	}
	afterOutText := string(parsedOutFile.Bytes())
	diffText, err := diffText(beforeOutText, afterOutText, 3)
	if len(diffText) > 0 {
		fmt.Printf("Diff for %v\n%v\n", toFilePath, diffText)
		plan.addFileUpdate(&FileUpdate{toFilePath, beforeOutText, afterOutText})
	}
	return &plan, nil
}

func findOrCreateLocalsBlock(parsedFile *hclwrite.File) *hclwrite.Block {
	found := parsedFile.Body().FirstMatchingBlock("locals", []string{})
	if found != nil {
		return found
	}
	return parsedFile.Body().AppendNewBlock("locals", []string{})
}

func moveLocals(parsedInFile, parsedOutFile *hclwrite.File) bool {
	for _, block := range parsedInFile.Body().Blocks() {

		if block.Type() == "locals" {
			fmt.Printf("## Found type locals\n")

			for attrKey, attrVal := range block.Body().Attributes() {
				fmt.Printf("adding %v\n", attrKey)
				toLocalsBlock := findOrCreateLocalsBlock(parsedOutFile)
				toLocalsBlock.Body().AppendUnstructuredTokens(attrVal.BuildTokens(nil))
			}
			if !parsedInFile.Body().RemoveBlock(block) {
				fmt.Printf("WARN locals block could not be removed\n")
			}
			return true
		}
	}
	return false
}

func moveLocal(localName string, parsedInFile, parsedOutFile *hclwrite.File) bool {
	for _, block := range parsedInFile.Body().Blocks() {

		if block.Type() == "locals" {
			attr := block.Body().GetAttribute(localName)
			if attr != nil {
				toLocalsBlock := findOrCreateLocalsBlock(parsedOutFile)
				toLocalsBlock.Body().AppendUnstructuredTokens(attr.BuildTokens(nil))

				block.Body().RemoveAttribute(localName)
				// This can leave an empty block. Maybe check for that.
				return true
			}

		}
	}
	return false
}

func labelsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func moveBlock(addr *Address, parsedInFile, parsedOutFile *hclwrite.File) bool {
	found := false
	addrLabels := addr.labels
	for _, block := range parsedInFile.Body().Blocks() {
		blockLabelsLimited := block.Labels()[0:min(len(addrLabels), len(block.Labels()))]
		if string(addr.BlockType()) == block.Type() && matchLabels(addr.labels, blockLabelsLimited) {
			// fmt.Printf("## Block matched %v %v\n", block.Type(), block.Labels())
			parsedOutFile.Body().AppendNewline()
			parsedOutFile.Body().AppendBlock(block)
			if !parsedInFile.Body().RemoveBlock(block) {
				fmt.Printf("WARN locals block could not be removed\n")
			}
			found = true
		}
	}
	return found
}

func moveAddrToFile(addr *Address, parsedInFile, parsedOutFile *hclwrite.File) bool {
	if addr.elementType == TypeLocal && len(addr.labels) == 0 {
		return moveLocals(parsedInFile, parsedOutFile)
	} else if addr.elementType == TypeLocal {
		localName := addr.labels[0]
		return moveLocal(localName, parsedInFile, parsedOutFile)
	}
	return moveBlock(addr, parsedInFile, parsedOutFile)
}
