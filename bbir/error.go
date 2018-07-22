package bbir

import (
	"strings"

	"github.com/golang/go/src/sort"
)

const (
	UnexpectedError = iota + 3
	BacklogAPIRequestError
	FileOpenError
	FileReadError
	ValidationIssueError
	RegistrationIssueError
)

func NewMultipleErrorsMap(errsMap map[string]error) *MultipleErrorsMap {
	return &MultipleErrorsMap{errsMap}
}

type MultipleErrorsMap struct {
	errsMap map[string]error
}

func (e MultipleErrorsMap) Error() string {
	var keys []string
	for k := range e.errsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString(":\n")
		b.WriteString(e.errsMap[k].Error())
	}
	return b.String()
}

func NewMultipleErrors(errs []error) *MultipleErrors {
	return &MultipleErrors{errs}
}

type MultipleErrors struct {
	errs []error
}

func (e MultipleErrors) Error() string {
	var b strings.Builder
	for _, err := range e.errs {
		b.WriteString("- ")
		b.WriteString(err.Error())
		b.WriteString("\n")
	}
	return b.String()
}
