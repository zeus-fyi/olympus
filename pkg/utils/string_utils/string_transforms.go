package string_utils

import (
	"net/url"
	"strings"
)

func DelimitedSliceStrBuilder(prefix, delimiter string, values ...string) string {
	var sb strings.Builder
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}
	for i, val := range values {
		sb.WriteString(val)

		if len(values)-1 == i {
			return sb.String()
		}
		sb.WriteString(delimiter)
	}
	return sb.String()
}

func UrlEncodeQueryParamList(prefix string, values ...string) string {
	return DelimitedSliceStrBuilder(prefix, "%2C", values...)
}

func UrlExplicitEncodeQueryParamList(key string, values ...string) string {
	u := url.Values{}
	for _, val := range values {
		u.Add(key, val)
	}
	return u.Encode()
}

func UrlPathStrBuilder(values ...string) string {
	return DelimitedSliceStrBuilder("", "/", values...)
}

func SliceStrBuilderCSV(values ...string) string {
	return DelimitedSliceStrBuilder("", ",", values...)
}

func SliceStrBuilderWithPrefixCSV(prefix string, values ...string) string {
	return DelimitedSliceStrBuilder(prefix, ",", values...)
}

func InsertValuesSliceStrBuilderWithSQLPrefix(prefix string, values ...string) string {
	return DelimitedSliceStrBuilder(prefix, ",", values...)
}
