package regex

import "regexp"

// Split slices s into substrings separated by the regex and returns a slice of substrings.
//
// The result contains all part of s, not just the part matched regex.
func Split(s string, re *regexp.Regexp) []string {
	var result []string
	index := 0
	indices := re.FindAllStringIndex(s, -1)
	for _, pair := range indices {
		if index < pair[0] {
			result = append(result, s[index:pair[0]])
		}
		result = append(result, s[pair[0]:pair[1]])
		index = pair[1]
	}
	if index < len(s) {
		result = append(result, s[index:])
	}
	return result
}
