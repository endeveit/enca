# enca [![Build Status](https://travis-ci.org/endeveit/enca.svg?branch=master)](https://travis-ci.org/endeveit/enca)

This is a minimal cgo bindings for [libenca](http://cihar.com/software/enca/).

If you need to detect the language of a string you can use [guesslanguage](https://github.com/endeveit/guesslanguage) package.

## Supported Go versions

go-enca is tested against Go 1.0, 1.1, 1.2, 1.3 and tip.

## Usage

Install libenca to your system:
```
$ sudo apt-get install libenca0 libenca-dev
```

Install in your `${GOPATH}` using `go get -u github.com/endeveit/enca`

Then call it:
```go
package main

import (
	"fmt"
	"github.com/endeveit/enca"
)

func main() {
	analyzer, err := enca.New("zh")

	if err == nil {
		encoding, err := analyzer.FromString("美国各州选民今天开始正式投票。据信，", enca.NAME_STYLE_HUMAN)
		defer analyzer.Free()

		// Output:
		// UTF-8
		if err == nil {
			fmt.Println(encoding)
		}
	}
}
```

## Documentation

godoc is available [here](http://godoc.org/github.com/endeveit/enca).
