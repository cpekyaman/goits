package validation

import (
	"regexp"
	"strings"
)

const (
	patternAlphaRegex            = "^[a-zA-Z]+$"
	patternAlnumRegex            = "^[a-zA-Z0-9]+$"
	patternDigitRegex            = "^[0-9]+$"
	patternWordAlphaRegex        = "^[a-zA-Z-_ ]+$"
	patternWordAlnumRegex        = "^[a-zA-Z0-9-_ ]+$"
	patternUnicodeAlphaRegex     = `^[\p{L}]+$`
	patternUnicodeWordAlphaRegex = `^[\p{L}-_ ]+$`
	patternUnicodeAlnumRegex     = `^[\p{L}\p{N}]+$`
	patternUnicodeWordAlnumRegex = `^[\p{L}\p{N}-_ ]+$`
)

// PatternType represents the enum of registered string validation patterns.
type PatternType uint8

const (
	PatternAlpha PatternType = iota
	PatternAlNum
	PatternDigit
	PatternWordAlpha
	PatternWordAlnum
	PatternUnicodeAlpha
	PatternUnicodeAlnum
	PatternUnicodeWordAlpha
	PatternUnicodeWordAlnum
)

var patternNames = [...]string{
	"Alpha", "AlphaNumeric", "Digit", "WordAlpha", "WordAlnum",
	"UnicodeAlpha", "UnicodeAlnum", "UnicodeWordAlpha", "UnicodeWordAlnum"}

// String returns the string label for this pattern type.
func (this PatternType) String() string {
	return patternNames[this]
}

var patternRegistry map[PatternType]*regexp.Regexp

func init() {
	patternRegistry = make(map[PatternType]*regexp.Regexp)

	patternRegistry[PatternAlpha] = regexp.MustCompile(patternAlphaRegex)
	patternRegistry[PatternAlNum] = regexp.MustCompile(patternAlnumRegex)
	patternRegistry[PatternDigit] = regexp.MustCompile(patternDigitRegex)
	patternRegistry[PatternWordAlpha] = regexp.MustCompile(patternWordAlphaRegex)
	patternRegistry[PatternWordAlnum] = regexp.MustCompile(patternWordAlnumRegex)
	patternRegistry[PatternUnicodeAlpha] = regexp.MustCompile(patternUnicodeAlphaRegex)
	patternRegistry[PatternUnicodeAlnum] = regexp.MustCompile(patternUnicodeAlnumRegex)
	patternRegistry[PatternUnicodeWordAlpha] = regexp.MustCompile(patternUnicodeWordAlphaRegex)
	patternRegistry[PatternUnicodeWordAlnum] = regexp.MustCompile(patternUnicodeWordAlnumRegex)
}

var notEmpty = validatorImpl{
	name: "notempty",
	vFunc: strPredicateFunc(func(str string) bool {
		return len(str) > 0
	}),
}

// NotEmpty gets the validator that validates if a string is not empty.
func NotEmpty() validatorImpl {
	return notEmpty
}

var notBlank = validatorImpl{
	name: "notblank",
	vFunc: strPredicateFunc(func(str string) bool {
		return len(strings.TrimSpace(str)) > 0
	}),
}

// NotBlank gets the validator that validates if a string is not blank.
func NotBlank() validatorImpl {
	return notBlank
}

// Pattern gets the validator that validates if a string matches the required pattern.
func Pattern(pt PatternType) validatorImpl {
	return validatorImpl{
		name:   "pattern",
		params: map[string]interface{}{"pattern": pt.String()},
		vFunc: strPredicateFunc(func(str string) bool {
			return patternRegistry[pt].Match([]byte(str))
		}),
	}
}

// StrLen gets the validator that validates if the length of a string is in given range.
func StrLen(min int, max int) validatorImpl {
	return validatorImpl{
		name:   "length",
		params: map[string]interface{}{"min": min, "max": max},
		vFunc: strPredicateFunc(func(str string) bool {
			l := len(str)

			if min > 0 && l < min {
				return false
			}
			if max > 0 && l > max {
				return false
			}
			return true
		}),
	}
}

func strPredicateFunc(pred func(string) bool) func(interface{}) bool {
	return func(value interface{}) bool {
		in := fieldValue(value)
		str, ok := in.(string)
		if !ok {
			return false
		}
		return pred(str)
	}
}
