package replace

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/apaxa-go/eval"
	"gopkg.in/yaml.v3"
)

// Needle search & replace options for a single operation.
type Needle struct {
	s       []byte `yaml:"-"`
	Search  string `yaml:"s,omitempty"`
	SCmd    string `yaml:"sc,omitempty"`
	SEval   string `yaml:"se,omitempty"`
	SRegex  string `yaml:"sx,omitempty"`
	r       []byte `yaml:"-"`
	Replace string `yaml:"r,omitempty"`
	RCmd    string `yaml:"rc,omitempty"`
	REval   string `yaml:"re,omitempty"`
}

func ReplaceYAMLFile(src []byte, path string) []byte {
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return ReplaceYAML(src, yml)
}

func ReplaceYAMLStr(src []byte, yml string) []byte {
	return ReplaceYAML(src, []byte(yml))
}

func ReplaceYAML(src, yml []byte) []byte {
	var h []Needle
	err := yaml.Unmarshal(yml, &h)
	if err != nil {
		log.Fatalln(err)
	}
	return Replace(src, h)
}

// Replace returns the modified byte slice src. Errors are printed to stderr & Replace continues processing.
func Replace(src []byte, hs []Needle) []byte {
	var searchRegex *regexp.Regexp
	var err error
	var isRegex bool
	for _, h := range hs {
		isRegex = false
		switch {
		case h.SRegex != "":
			isRegex = true
			searchRegex, err = regexp.Compile(h.SRegex)
		case h.SCmd != "":
			cmds := strings.Split(h.SCmd, " ")
			h.s, err = run(cmds[0], cmds[1:]...)
		case h.SEval != "":
			h.s, err = evl(h.SEval)
		default:
			h.s = []byte(h.Search)
			err = nil
		}
		if err != nil {
			log.Println(err)
			continue
		}

		switch {
		case h.RCmd != "":
			cmds := strings.Split(h.RCmd, " ")
			h.r, err = run(cmds[0], cmds[1:]...)
		case h.REval != "":
			h.r, err = evl(h.REval)
		default:
			h.r = []byte(h.Replace)
			err = nil
		}
		if err != nil {
			log.Println(err)
			continue
		}

		if isRegex {
			src = searchRegex.ReplaceAllLiteral(src, h.r)
		} else {
			src = bytes.Replace(src, h.s, h.r, -1)
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
		log.Println("Error with eval:", src, ":", err)
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
