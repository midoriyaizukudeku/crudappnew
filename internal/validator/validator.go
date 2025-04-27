package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldError []string
	FieldsErrors  map[string]string
}

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z]{2,})+$`)

func (v *Validator) Valid() bool {
	return len(v.FieldsErrors) == 0 && len(v.NonFieldError) == 0
}
func (v *Validator) addNonfoelderror(message string) {
	v.NonFieldError = append(v.NonFieldError, message)
}
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldsErrors == nil {
		v.FieldsErrors = make(map[string]string)
	}
	if _, exists := v.FieldsErrors[key]; !exists {
		v.FieldsErrors[key] = message
	}

}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChar(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValues[T comparable](value T, permittedvalues ...T) bool {
	for i := range permittedvalues {
		if value == permittedvalues[i] {
			return true
		}
	}
	return false
}

func MinChar(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
