package enca

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_GetAvailableLanguages(t *testing.T) {
	if len(GetAvailableLanguages()) == 0 {
		t.Error("List of available languages must not be empty")
	}
}

func Test_EncaFromString(t *testing.T) {
	analyzer, err := New("zh")

	if err != nil {
		t.Errorf("Unable to create new analyzer: %v", err)
	} else {
		encoding, err := analyzer.FromString("美国各州选民今天开始正式投票。据信，", NAME_STYLE_HUMAN)
		defer analyzer.Free()

		// Output:
		// UTF-8
		if err != nil {
			t.Errorf("Unable to get encoding: %v", err)
		} else {
			if encoding != "UTF-8" {
				t.Errorf("Encoding must be 'UTF-8'")
			}
		}
	}
}

func Test_EncaFromBytes(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Unable to get current working directory: %v", err)
	}

	var (
		ea      *EncaAnalyser
		parts   []string
		lang    string
		tests   map[string][]string = make(map[string][]string)
		content []byte
	)

	visit := func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			parts = strings.Split(f.Name(), "-")
			if len(parts) < 2 {
				return fmt.Errorf("Unable to process file '%s'", path)
			}

			lang = parts[0]
			if lang == "none" {
				lang = "__"
			}

			if _, ok := tests[lang]; ok {
				tests[lang] = append(tests[lang], path)
			} else {
				tests[lang] = []string{path}
			}
		}

		return err
	}

	err = filepath.Walk(fmt.Sprintf("%s%stest-data", pwd, string(filepath.Separator)), visit)
	if err != nil {
		t.Error(err)
	}

	if len(tests) == 0 {
		t.Errorf("Tests are empty")
	}

	for lang, files := range tests {
		ea, err = New(lang)
		if err != nil {
			t.Errorf("Unable to create new analyzer: %v", err)
			break
		}

		for _, file := range files {
			content, err = ioutil.ReadFile(file)
			if err != nil {
				t.Errorf("Unable to read file '%s': %v", file, err)
			} else {
				_, err = ea.FromBytes(content, NAME_STYLE_ENCA)
				if err != nil {
					t.Errorf("Unable to detect encoding of '%s': %v", file, err)
				}
			}

		}

		ea.Free()
	}
}
