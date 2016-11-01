package util

import "regexp"

// From https://github.com/peco/peco/blob/d8e5bcb/internal/util/util.go#L22
// Global var used to strips ansi sequences
var reANSIEscapeChars = regexp.MustCompile("\x1B\\[(?:[0-9]{1,2}(?:;[0-9]{1,2})?)*[a-zA-Z]")

// StripANSISequence strips ansi sequences
func StripANSISequence(s string) string {
	return reANSIEscapeChars.ReplaceAllString(s, "")
}
