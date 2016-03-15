// Package enca provides a minimal cgo bindings for libenca
//
// Source code and project home:
// https://github.com/endeveit/go-enca
//
package enca

/*
#cgo pkg-config: enca

#include <enca.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

const (
	// Default, implicit charset name in Enca
	NAME_STYLE_ENCA = 0 << iota
	// RFC 1345 or otherwise canonical charset name
	NAME_STYLE_RFC1345
	// Cstocs charset name (may not exist)
	NAME_STYLE_CSTOCS
	// Iconv charset name (may not exist)
	NAME_STYLE_ICONV
	// Human comprehensible description
	NAME_STYLE_HUMAN
	// Preferred MIME name (may not exist)
	NAME_STYLE_MIME
)

var availableLanguages []string

func init() {
	n := C.size_t(0)
	l := C.enca_get_languages(&n)

	availableLanguages = getStrings(l, n)

	C.free(unsafe.Pointer(l))
}

type EncaAnalyser struct {
	sync.Mutex
	// Language for which the analyser is initialized
	Language string
	enca     C.EncaAnalyser
}

// Returns list of available languages
func GetAvailableLanguages() []string {
	return availableLanguages
}

// Returns a new EncaAnalyzer object for the given language.
func New(lang string) (*EncaAnalyser, error) {
	if !keyExists(lang, availableLanguages) {
		return nil, fmt.Errorf(
			"Invalid language '%s'. Available languages are: '%s'",
			lang,
			strings.Join(availableLanguages, "', '"))
	}

	cLang := C.CString(lang)
	cOne := C.int(1)

	defer C.free(unsafe.Pointer(cLang))

	analyzer := &EncaAnalyser{Language: lang, enca: C.enca_analyser_alloc(cLang)}
	C.enca_set_threshold(analyzer.enca, C.double(1.38))
	C.enca_set_multibyte(analyzer.enca, cOne)
	C.enca_set_ambiguity(analyzer.enca, cOne)
	C.enca_set_garbage_test(analyzer.enca, cOne)

	return analyzer, nil
}

// Returns encoding of provided byte array
func (ea *EncaAnalyser) FromBytes(bytes []byte, nameStyle int) (result string, err error) {
	ea.Lock()
	defer ea.Unlock()
	cText := (*C.uchar)(unsafe.Pointer(&bytes[0]))

	encoding := C.enca_analyse_const(ea.enca, cText, C.size_t(len(bytes)))

	if encoding.charset == C.int(-1) {
		errno := C.enca_errno(ea.enca)
		err = fmt.Errorf("Error %d: %s", errno, C.GoString(C.enca_strerror(ea.enca, errno)))
	} else {
		result = C.GoString(C.enca_charset_name(encoding.charset, C.EncaNameStyle(nameStyle)))
	}

	return result, err
}

// Helper function that returns encoding of provided string
func (ea *EncaAnalyser) FromString(text string, nameStyle int) (string, error) {
	return ea.FromBytes([]byte(text), nameStyle)
}

// Frees memory used by EncaAnalyser
func (ea *EncaAnalyser) Free() {
	C.enca_analyser_free(ea.enca)
}

// Helper function that converts C "const char**" into Go []string
func getStrings(x **C.char, n C.size_t) []string {
	q := ((*[1 << 30]*C.char)(unsafe.Pointer(x)))[:n]
	r := make([]string, n)

	for i, cs := range q {
		r[i] = C.GoString(cs)
	}

	return r
}

// Check if key exists
func keyExists(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
