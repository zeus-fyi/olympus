package utils

import (
	"strings"

	"github.com/zeus-fyi/olympus/databases/postgres"
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

func DelimitedInterfaceSliceStrBuilderSQL(sb *strings.Builder, delimiter string, values postgres.RowValues) *strings.Builder {
	for i, val := range values {
		sb.WriteString(val)

		if len(values)-1 == i {
			return sb
		}
		sb.WriteString(delimiter)
	}
	return sb
}

func SQLDelimitedSliceStrBuilder(prefix string, entries postgres.RowEntries) string {
	var sb strings.Builder
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}

	for count, row := range entries.Rows {

		sb.WriteString("(")
		DelimitedInterfaceSliceStrBuilderSQL(&sb, ",", row)
		sb.WriteString(")")

		if len(entries.Rows)-1 == count {
			return sb.String()
		}
		sb.WriteString(",")
	}

	return sb.String()
}

func SliceStrBuilderWithPrefixCSV(prefix string, values ...string) string {
	return DelimitedSliceStrBuilder(prefix, ",", values...)
}

func InsertValuesSliceStrBuilderWithSQLPrefix(prefix string, values ...string) string {
	return DelimitedSliceStrBuilder(prefix, ",", values...)
}

func UrlEncodeQueryParamList(prefix string, values ...string) string {
	return DelimitedSliceStrBuilder(prefix, "%2C", values...)
}

func SliceStrBuilderCSV(values ...string) string {
	return DelimitedSliceStrBuilder("", ",", values...)
}

func UrlPathStrBuilder(values ...string) string {
	return DelimitedSliceStrBuilder("", "/", values...)
}
