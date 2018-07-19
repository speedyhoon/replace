package replace

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"path/filepath"
)

func TestReplace(t *testing.T) {
	input, err := ioutil.ReadFile(filepath.Join("test", "input.html"))
	if err != nil {
		t.Error("unable to read input.html")
	}
	expected, err := ioutil.ReadFile(filepath.Join("test", "expected.html"))
	if err != nil {
		t.Error("unable to read expected.html")
	}

	settings := []Needle{
		{Search: "{{form .", Replace: "{{template ."},
		{SRegex: "<(/){0,1}samp>", Replace: "<s>"},
		{Search: "go help tool", RCmd: "go help tool"},
		{Search: ".<", REval: "int8( 3 * (1 + 2) )"},
	}

	output := Replace(input, settings)

	fmt.Println(len(expected), len(output))

	if string(expected) == string(output) {
		println("strings are different too")
	} else {
		println("strings are the same!")
	}

	fmt.Printf("Expected: %s\n", expected)
	fmt.Printf("Output: %s\n", output)

	if !bytes.Equal(expected, output) {
		t.Error("Output differs from expected\n", printDiff(expected, output))
	}
}

func printDiff(original, compareTo []byte) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(compareTo), string(original), false)
	if len(diffs) > 1 {
		return dmp.DiffPrettyText(diffs)
	}
	return ""
}
