package string_utils

import "strings"

func FilterString(word, contains string, doesNotContain []string) bool {
	passFilter := false

	// empty doesNotContain in means ignore, also must have at least one letter
	for _, check := range doesNotContain {
		if len(check) > 0 && strings.Contains(word, check) {
			return false
		}
	}

	if strings.Contains(word, contains) {
		passFilter = true
	}
	return passFilter
}
