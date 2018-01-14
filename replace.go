package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/apaxa-go/eval"
	"gopkg.in/yaml.v2"
)

var wrn = log.New(os.Stderr, "", 0)

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println(`
List of search and replace options:
	s:  search string
	sc: run command & search for the returned output
	se: evaluate Go code & search for the returned output
	sx: search regex
	r:  replace string
	rc: run command & replace with the returned output
	re: evaluate Go code & replace with the returned output
`)
	}

	file := flag.String("file", "", `Load the specified YAML file with a list of search & replace options
	Example file contents to replace "foo" with "bar" and "9" with "nine":
	- s: foo,
	  r: bar
	- se: int8((1+2)*3),
	  r: nine
`)
	yml := flag.String("yaml", "", `Example inline parameters to replace "foo" with "bar" and "9" with "nine":
	-yaml="[{s: foo, r: bar},{se: int8((1+2)*3), r: nine}]"`)
	flag.Parse()

	var h []Haystack
	err := yaml.Unmarshal([]byte(*yml), &h)
	if err != nil {
		wrn.Fatalln(err)
	}

	if len(*file) > 0 {
		yamlSrc, err := ioutil.ReadFile(*file)
		if err != nil {
			wrn.Fatalln(err)
		}

		var j []Haystack
		err = yaml.Unmarshal(yamlSrc, &j)
		if err != nil {
			wrn.Fatalln(err)
		}

		h = append(h, j...)
	}

	//Read Standard In stream
	src, err := bufio.NewReader(os.Stdin).ReadBytes(0)
	if err != nil && err != io.EOF {
		wrn.Fatalln(err)
	}

	os.Stdout.Write(Replace(src, h))
}

//Haystack search & replace options for a single operation.
type Haystack struct {
	search      []byte `yaml:"-"`
	Search      string `yaml:"s,omitempty"`
	SearchCmd   string `yaml:"sc,omitempty"`
	SearchEval  string `yaml:"se,omitempty"`
	SearchRegex string `yaml:"sx,omitempty"`
	replace     []byte `yaml:"-"`
	Replace     string `yaml:"r,omitempty"`
	ReplaceCmd  string `yaml:"rc,omitempty"`
	ReplaceEval string `yaml:"re,omitempty"`
}

//Replace returns the modified byte slice src. If any Haystack options error they are printed to stderr & Replace continues processing.
func Replace(src []byte, hs []Haystack) []byte {
	var searchRegex *regexp.Regexp
	var err error
	var isRegex bool
	for _, h := range hs {
		isRegex = false
		switch {
		case h.SearchRegex != "":
			isRegex = true
			searchRegex = regexp.MustCompile(h.SearchRegex)
		case h.SearchCmd != "":
			s := strings.Split(h.SearchCmd, " ")
			h.search, err = run(s[0], s[1:]...)
			if err != nil {
				wrn.Println("Error with cmd", h.SearchCmd, ":", err)
				continue
			}
		case h.SearchEval != "":
			h.search, err = evl(h.SearchEval)
			if err != nil {
				continue
			}
		default:
			h.search = []byte(h.Search)
		}

		switch {
		case h.ReplaceCmd != "":
			s := strings.Split(h.ReplaceCmd, " ")
			h.replace, err = run(s[0], s[1:]...)
			if err != nil {
				wrn.Println("Error with", h.ReplaceCmd, ":", err)
				continue
			}
		case h.ReplaceEval != "":
			h.replace, err = evl(h.ReplaceEval)
			if err != nil {
				continue
			}
		default:
			h.replace = []byte(h.Replace)
		}

		if isRegex {
			src = searchRegex.ReplaceAllLiteral(src, h.replace)
		} else {
			src = bytes.Replace(src, h.search, h.replace, -1)
		}
	}
	return src
}

func evl(src string) ([]byte, error) {
	expr, err := eval.ParseString(src, "")
	if err != nil {
		return nil, err
	}
	r, err := expr.EvalToInterface(nil)
	if err != nil {
		wrn.Println("Error with eval:", src, ":", err)
	}
	return []byte(fmt.Sprintf("%v", r)), err
}

func run(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	var sOut, sErr bytes.Buffer
	cmd.Stdout = &sOut
	cmd.Stderr = &sErr

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	if sErr.String() != "" {
		return sOut.Bytes(), errors.New(sErr.String())
	}

	return sOut.Bytes(), nil
}
