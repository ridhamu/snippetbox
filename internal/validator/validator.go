// Package validator contains validation tools
package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddField(key string, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, ok := v.FieldErrors[key]; !ok {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddField(key, message)
	}
}

// NotBlankString spesific validation perform here
func NotBlankString(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinChars check wether given values is >= given min
func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

// MaxChars check wehter given values is <= given max
func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// PermittedValue check if given value is one of the given permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
