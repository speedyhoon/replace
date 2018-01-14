package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestReplace(t *testing.T) {

	input, err := ioutil.ReadFile("input.html")
	if err != nil {
		t.Error("unable to read input.html")
	}
	expected, err := ioutil.ReadFile("expected.html")
	if err != nil {
		t.Error("unable to read expected.html")
	}

	settings := []Haystack{
		{Search: "{{form .", Replace: "{{template ."},
		{SearchRegex: "<(/){0,1}samp>", Replace: "<s>"},
		{Search: "go help tool", ReplaceCmd: "go help tool"},
		{Search: ".<", ReplaceEval: "int8( 3 * (1 + 2) )"},
	}

	output := Replace(input, settings)

	fmt.Println(len(expected), len(output))

	if string(expected) == string(output) {
		println("strings are different too")
	}else{
		println("strings are the same!")
	}

/*
	for i := 0; i < len(expected); i++{
		if expected[i] != output[i]{
			print("different ", i, expected[i], output[i])
			fmt.Printf("%s == %s\n", expected[i], output[i])
		}
	}*/

	fmt.Printf("Expected: %s\n", expected)
	fmt.Printf("Output: %s\n", output)

	//ioutil.WriteFile("temp.html", output, 0777)
	//if err != nil {
	//	t.Error("unable to write temp.html")
	//}

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
