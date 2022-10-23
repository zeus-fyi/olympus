package string_utils

import (
	"strings"
)

func FilterStringWithOpts(word string, filter *FilterOpts) bool {
	// empty doesNotContain in means ignore, also must have at least one letter
	if filter == nil {
		return true
	}
	for _, check := range filter.DoesNotInclude {
		if len(check) > 0 && strings.Contains(word, check) {
			return false
		}
	}

	if len(filter.StartsWith) > 0 {
		if !strings.HasPrefix(word, filter.StartsWith) {
			return false
		}
	}

	if len(filter.Contains) > 0 {
		if !strings.Contains(word, filter.Contains) {
			return false
		}
	}

	if len(filter.DoesNotStartWith) > 0 {
		for _, wordFilter := range filter.DoesNotStartWith {
			if strings.HasPrefix(word, wordFilter) {
				return false
			}
		}
	}
	return true
}

type FilterOpts struct {
	DoesNotStartWith []string
	StartsWith       string
	Contains         string
	DoesNotInclude   []string
}
