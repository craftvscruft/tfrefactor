package refactor

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMv(t *testing.T) {
	cases := []struct {
		name string
		args []string
		ok   bool
		from string
		to   string
	}{
		{
			name: "case_data_single_block",
			args: []string{"data", "data.tf"},
			ok:   true,
			from: "test_data/mv/case_data_single_block/from",
			to:   "test_data/mv/case_data_single_block/to",
		},
		{
			name: "case_data_single_block_qualified",
			args: []string{"data.a", "data.tf"},
			ok:   true,
			from: "test_data/mv/case_data_single_block/from",
			to:   "test_data/mv/case_data_single_block/to",
		},
		{
			name: "case_data_single_block_fully_qualified",
			args: []string{"data.a.b", "data.tf"},
			ok:   true,
			from: "test_data/mv/case_data_single_block/from",
			to:   "test_data/mv/case_data_single_block/to",
		},
		{
			name: "case_data_single_block_new_file",
			args: []string{"data", "data.tf"},
			ok:   true,
			from: "test_data/mv/case_data_single_block_new_file/from",
			to:   "test_data/mv/case_data_single_block_new_file/to",
		},
		{
			name: "case_all_data_blocks",
			args: []string{"data", "data.tf"},
			ok:   true,
			from: "test_data/mv/case_all_data_blocks/from",
			to:   "test_data/mv/case_all_data_blocks/to",
		},
		{
			name: "case_all_vars",
			args: []string{"variable", "variables.tf"},
			ok:   true,
			from: "test_data/mv/case_all_vars/from",
			to:   "test_data/mv/case_all_vars/to",
		},
		{
			name: "case_single_var",
			args: []string{"variable.a", "variables.tf"},
			ok:   true,
			from: "test_data/mv/case_single_var/from",
			to:   "test_data/mv/case_single_var/to",
		},
		{
			name: "case_all_locals",
			args: []string{"locals", "locals.tf"},
			ok:   true,
			from: "test_data/mv/case_all_locals/from",
			to:   "test_data/mv/case_all_locals/to",
		},
		{
			name: "case_single_local",
			args: []string{"local.b", "locals.tf"},
			ok:   true,
			from: "test_data/mv/case_single_local/from",
			to:   "test_data/mv/case_single_local/to",
		},
		{
			name: "case_resource_type",
			args: []string{"resource.a", "dest.tf"},
			ok:   true,
			from: "test_data/mv/case_resource_type/from",
			to:   "test_data/mv/case_resource_type/to",
		},
		{
			name: "case_all_outputs",
			args: []string{"output", "outputs.tf"},
			ok:   true,
			from: "test_data/mv/case_all_outputs/from",
			to:   "test_data/mv/case_all_outputs/to",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := Mv(tc.args[0], tc.args[1], tc.from)
			if (err == nil) != tc.ok {

			}
			assertMockCmdFileOutput(t, tc.name, tc.from, tc.to, plan)
		})
	}
}

func copyFile(src string, dest string) {
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		log.Fatal(err)
	}
}

func fileNames(vs []os.FileInfo) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = v.Name()
	}
	return vsm
}

// assertMockCmd is a high-level test helper to run a given mock command with
// arguments and check if an error and its stdout are expected.
func assertMockCmdFileOutput(t *testing.T, name string, from string, to string, plan *UpdatePlan) {
	expectedFiles, err := ioutil.ReadDir(to)
	if err != nil {
		log.Fatal(err)
	}
	startingFiles, err := ioutil.ReadDir(from)

	filenameToContents := make(map[string]string)
	for _, file := range startingFiles {
		filePath := filepath.Join(from, file.Name())
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Loading %v - %v", filePath, err)
		}
		filenameToContents[filepath.Base(file.Name())] = string(buf)
	}
	for _, update := range plan.FileUpdates {
		filenameToContents[filepath.Base(update.Filename)] = update.AfterText
	}

	if len(expectedFiles) != len(filenameToContents) {
		actualFilenames := []string{}
		for k := range filenameToContents {
			actualFilenames = append(actualFilenames, k)
		}
		t.Fatalf("Expected files to be %v, but found %v", fileNames(expectedFiles), actualFilenames)
	}
	for _, file := range expectedFiles {
		assertFileHasContents(t, to+"/"+file.Name(), filenameToContents[filepath.Base(file.Name())])
	}
}

func assertFileHasContents(t *testing.T, expectedFile string, actualContents string) {
	expectedBuf, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		log.Fatalf("Loading %v - %v", expectedFile, err)
	}
	expected := strings.TrimSpace(string(normalizeNewlines(expectedBuf)))
	actual := strings.TrimSpace(string(normalizeNewlines([]byte(actualContents))))
	assert.Equal(t, expected, actual, "File %v", expectedFile)
}

// NormalizeNewlines normalizes \r\n (windows) and \r (mac)
// into \n (unix)
func normalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}
