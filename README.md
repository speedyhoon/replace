# Replace
[![Build Status](https://travis-ci.org/speedyhoon/replace.svg?branch=master)](https://travis-ci.org/speedyhoon/replace)
[![go report card](https://goreportcard.com/badge/github.com/speedyhoon/replace)](https://goreportcard.com/report/github.com/speedyhoon/replace)

Cmdline search replace tool files or stdin &amp; stdout

cat file.txt | replace -yaml="[{s: cats, r: dogs}]"
cat file.txt | replace -json={"s":"cats","r":"dogs"}