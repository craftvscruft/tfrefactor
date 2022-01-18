package refactor

// Adapted from hcledit by Masayuki Morita

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/raymyers/hcl/v2"
	"github.com/raymyers/hcl/v2/hclwrite"
)

func ParseHclFile(filename string) (*hclwrite.File, error) {
	input, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return safeParseConfig(input, filename, hcl.Pos{Line: 1, Column: 1})
}

// Source parses HCL and returns *hclwrite.File
// filename is a metadata of input stream and used only for an error message.
func ParseHclBytes(src []byte, filename string) (*hclwrite.File, error) {
	return safeParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
}

// safeParseConfig parses config and recovers if panic occurs.
// The current hclwrite implementation is no perfect and will panic if
// unparseable input is given. We just treat it as a parse error so as not to
// surprise users.
func safeParseConfig(src []byte, filename string, start hcl.Pos) (f *hclwrite.File, e error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[DEBUG] failed to parse input: %s\nstacktrace: %s", filename, string(debug.Stack()))
			// Set a return value from panic recover
			e = fmt.Errorf(`failed to parse input: %s
panic: %s
This may be caused by a bug in the hclwrite parser`, filename, err)
		}
	}()

	f, diags := hclwrite.ParseConfig(src, filename, start)

	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse input: %s", diags)
	}

	return f, nil
}
