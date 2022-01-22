package refactor

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	// from, err := filepath.Abs(from)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	expectedFiles := []string{}
	err := filepath.Walk(to,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				rel := relativise(to, path)
				expectedFiles = append(expectedFiles, rel)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	// expectedFiles, err := ioutil.ReadDir(to)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	startingFiles, err := ioutil.ReadDir(from)

	filenameToContents := make(map[string]string)
	for _, file := range startingFiles {
		filePath := filepath.Join(from, file.Name())
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Loading %v - %v", filePath, err)
		}
		filenameToContents[file.Name()] = string(buf)
	}
	for _, update := range plan.FileUpdates {
		filenameToContents[relativise(from, update.Filename)] = update.AfterText
	}

	if len(expectedFiles) != len(filenameToContents) {
		actualFilenames := []string{}
		for k := range filenameToContents {
			actualFilenames = append(actualFilenames, k)
		}
		t.Fatalf("Expected files to be %v, but found %v", expectedFiles, actualFilenames)
	}
	for _, file := range expectedFiles {
		if actualContents, ok := filenameToContents[file]; ok {
			assertFileHasContents(t, filepath.Join(to, file), actualContents)
		} else {
			t.Fatalf("Didn't find found %v in:\n%v", file, filenameToContents)
		}
	}
}

func relativise(from, name string) string {
	absUpdatePath, _ := filepath.Abs(name)
	absFromPath, _ := filepath.Abs(from)
	relUpdatePath, _ := filepath.Rel(absFromPath, absUpdatePath)
	return relUpdatePath
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
