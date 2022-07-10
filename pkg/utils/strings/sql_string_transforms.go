package strings

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/databases/postgres"
)

func StringDelimitedSliceBuilderSQL(sb *strings.Builder, delimiter string, values postgres.RowValues) *strings.Builder {
	for i, val := range values {

		switch val.(type) {
		case string:
			sb.WriteString("'")
			sb.WriteString(val.(string))
			sb.WriteString("'")
		case int, int64:
			sb.WriteString(fmt.Sprintf("%d", val.(int64)))
		case uint64:
			sb.WriteString(fmt.Sprintf("%d", val.(uint64)))
		case bool:
			sb.WriteString(fmt.Sprintf("%t", val.(bool)))
		default:
		}

		if len(values)-1 == i {
			return sb
		}
		sb.WriteString(delimiter)
	}
	return sb
}

func PrefixAndSuffixDelimitedSliceStrBuilderSQLRows(prefix string, entries postgres.RowEntries, suffix string) string {
	var sb strings.Builder
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}

	for count, row := range entries.Rows {

		sb.WriteString("(")
		StringDelimitedSliceBuilderSQL(&sb, ",", row)
		sb.WriteString(")")

		if len(entries.Rows)-1 == count {
			return sb.String()
		}
		sb.WriteString(",")
	}

	if len(suffix) > 0 {
		sb.WriteString(suffix)
	}

	return sb.String()
}

func DelimitedSliceStrBuilderSQLRows(prefix string, entries postgres.RowEntries) string {
	var sb strings.Builder
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}

	for count, row := range entries.Rows {

		sb.WriteString("(")
		StringDelimitedSliceBuilderSQL(&sb, ",", row)
		sb.WriteString(")")

		if len(entries.Rows)-1 == count {
			return sb.String()
		}
		sb.WriteString(",")
	}

	return sb.String()
}

func AnyArraySliceStrBuilderSQL(entries postgres.RowValues) string {
	var sb strings.Builder

	sb.WriteString("ANY(ARRAY[")
	StringDelimitedSliceBuilderSQL(&sb, ",", entries)
	sb.WriteString("])")

	return sb.String()
}

func ArraySliceStrBuilderSQL(entries postgres.RowValues) string {
	var sb strings.Builder

	sb.WriteString("ARRAY[")
	StringDelimitedSliceBuilderSQL(&sb, ",", entries)
	sb.WriteString("]")

	return sb.String()
}

func MultiArraySliceStrBuilderSQL(r postgres.RowEntries) string {
	var sb strings.Builder

	for count, row := range r.Rows {
		sb.WriteString("ARRAY[")
		StringDelimitedSliceBuilderSQL(&sb, ",", row)
		sb.WriteString("]")

		if len(r.Rows)-1 == count {
			return sb.String()
		}
		sb.WriteString(",")
	}

	return sb.String()
}
