package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/speedyhoon/replace"
	"gopkg.in/yaml.v3"
)

var wrn = log.New(os.Stderr, "", 0)

func main() {
	file := flag.String("file", "", `Load a YAML file containing search & replace options
	Example file contents to replace "foo" with "bar" and "9" with "nine":
	- s: foo,
	  r: bar
	- se: int8((1+2)*3),
	  r: nine
`)
	yml := flag.String("yaml", "", `Example inline parameters to replace "foo" with "bar" and "9" with "nine":
	-yaml="[{s: foo, r: bar},{se: int8((1+2)*3), r: nine}]"`)

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
	re: evaluate Go code & replace with the returned output`)
	}
	flag.Parse()

	var h []replace.Needle
	err := yaml.Unmarshal([]byte(*yml), &h)
	if err != nil {
		wrn.Fatalln(err)
	}

	if len(*file) > 0 {
		yamlSrc, err := ioutil.ReadFile(*file)
		if err != nil {
			wrn.Fatalln(err)
		}

		var j []replace.Needle
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

	os.Stdout.Write(replace.Replace(src, h))
}
