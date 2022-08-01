package string_utils

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
)

func StringDelimitedSliceBuilderSQL(delimiter string, values postgres.RowValues) string {
	returnStr := ""
	for i, val := range values {

		switch val.(type) {
		case string:
			returnStr += "'"
			returnStr += val.(string)
			returnStr += "'"
		case int, int64:
			returnStr += fmt.Sprintf("%d", val.(int64))

		case uint64:
			returnStr += fmt.Sprintf("%d", val.(uint64))
		case bool:
			returnStr += fmt.Sprintf("%t", val.(bool))
		default:
		}

		if len(values)-1 == i {
			return returnStr
		}
		returnStr += delimiter
	}
	return returnStr
}

func PrefixAndSuffixDelimitedSliceStrBuilderSQLRows(prefix string, entries postgres.RowEntries, suffix string) string {
	sb := strings.Builder{}
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}
	returnStr := prefix

	for count, row := range entries.Rows {

		sb.WriteString("(")
		sb.WriteString(StringDelimitedSliceBuilderSQL(",", row))
		sb.WriteString(")")
		returnStr += sb.String()

		if len(entries.Rows)-1 == count {
			return returnStr
		}
		sb.WriteString(",")
	}

	if len(suffix) > 0 {
		sb.WriteString(suffix)
	}
	returnStr += sb.String()
	return returnStr
}

func DelimitedSliceStrBuilderSQLRows(prefix string, entries postgres.RowEntries) string {
	sb := strings.Builder{}
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}

	for count, row := range entries.Rows {

		sb.WriteString("(")
		tmp := StringDelimitedSliceBuilderSQL(",", row)
		sb.WriteString(tmp)
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
	sb.WriteString(StringDelimitedSliceBuilderSQL(",", entries))

	sb.WriteString("])")

	return sb.String()
}

func ArraySliceStrBuilderSQL(entries postgres.RowValues) string {
	var sb strings.Builder

	sb.WriteString("ARRAY[")
	sb.WriteString(StringDelimitedSliceBuilderSQL(",", entries))
	sb.WriteString("]")

	return sb.String()
}

func MultiArraySliceStrBuilderSQL(r postgres.RowEntries) string {
	var sb strings.Builder

	for count, row := range r.Rows {
		sb.WriteString("ARRAY[")
		sb.WriteString(StringDelimitedSliceBuilderSQL(",", row))

		sb.WriteString("]")

		if len(r.Rows)-1 == count {
			return sb.String()
		}
		sb.WriteString(",")
	}

	return sb.String()
}
