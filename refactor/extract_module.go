package refactor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/raymyers/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func Extract(addresses []string, toFolder, configPath string) (*UpdatePlan, error) {
	filenames, err := filepath.Glob(configPath + "/*.tf")
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	plan := newUpdatePlan()
	var parsedOutFile *hclwrite.File
	toFile := filepath.Join(configPath, toFolder, "main.tf")
	if _, err := os.Stat(toFile); errors.Is(err, os.ErrNotExist) {
		parsedOutFile, err = ParseHclBytes([]byte{}, toFile)
		if err != nil {
			return nil, err
		}
	} else {
		parsedOutFile, err = ParseHclFile(toFile)
		if err != nil {
			return nil, err
		}
	}

	beforeOutText := string(parsedOutFile.Bytes())
	for _, filename := range filenames {
		fromPath, _ := filepath.Abs(filename)
		toPath, _ := filepath.Abs(toFile)
		if fromPath != "" && fromPath != toPath {
			parsedInFile, err := ParseHclFile(filename)
			beforeText := string(parsedInFile.Bytes())
			for _, fromAddressText := range addresses {
				fromAddress := ParseAddress(fromAddressText)

				if err != nil {
					return nil, err
				}

				if err != nil {
					return nil, err
				}
				if moveAddrToFile(fromAddress, parsedInFile, parsedOutFile) {
					if fromAddress.elementType == TypeResource {
						moduleName := filepath.Base(toFolder)
						toAddress := ParseAddress("module." + moduleName + "." + strings.Join(fromAddress.labels, "."))
						AddModuleBlock(parsedInFile, moduleName, toFolder)
						AddMovedBlock(parsedInFile, fromAddress, toAddress)
					}
				}
			}
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
		fmt.Printf("Diff for %v\n%v\n", toFile, diffText)
		plan.addFileUpdate(&FileUpdate{toFile, beforeOutText, afterOutText})
	}
	return &plan, nil
}

func AddModuleBlock(file *hclwrite.File, moduleName, toFolder string) {
	file.Body().AppendNewline()
	movedBlock := file.Body().AppendNewBlock("module", []string{moduleName})

	movedBlock.Body().SetAttributeValue("source", cty.StringVal(toFolder))
}
