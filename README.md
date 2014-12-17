# go-enca

Minimal cgo bindings for libenca.

If you need to detect language of string you can use [guesslanguage](https://github.com/endeveit/guesslanguage).

## Supported Go versions

go-enca is tested against Go 1.0, 1.1, 1.2, 1.3 and tip.

## Usage

Install libenca to your system:
```
$ sudo apt-get install libenca0 libenca-dev
```

Install in your `${GOPATH}` using `go get -u github.com/endeveit/go-enca`

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
		encoding, err := analyzer.FromString("美国各州选民今天开始正式投票。据信，", NAME_STYLE_HUMAN)
		defer analyzer.Free()

		// Output:
		// UTF-8
		if err == nil {
			fmt.Println(encoding)
		}
	}
}
```
