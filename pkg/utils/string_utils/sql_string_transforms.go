package string_utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func StringDelimitedSliceBuilderSQLMultiRowValues(delimiter string, values []apps.RowValues) string {
	sb := strings.Builder{}
	for count, row := range values {
		sb.WriteString("(")
		for i, val := range row {
			switch val.(type) {
			case string:
				strValue := val.(string)
				sb.WriteString("'")
				sb.WriteString(val.(string))
				sb.WriteString("'")
				if strValue == "{}" {
					sb.WriteString("::jsonb")
				}
			case int:
				returnStr := fmt.Sprintf("%d", val.(int))
				sb.WriteString(returnStr)
			case int64:
				returnStr := fmt.Sprintf("%d", val.(int64))
				sb.WriteString(returnStr)
			case uint64:
				returnStr := fmt.Sprintf("%d", val.(uint64))
				sb.WriteString(returnStr)
			case bool:
				returnStr := fmt.Sprintf("%t", val.(bool))
				sb.WriteString(returnStr)
			case time.Time:
				sb.WriteString("'(NOW())'")
			default:
			}
			if len(row)-1 != i {
				sb.WriteString(delimiter)
			}
		}

		sb.WriteString(")")
		if len(values)-1 == count {
			return sb.String()
		}
		sb.WriteString(delimiter)
	}
	return sb.String()
}

func StringDelimitedSliceBuilderSQL(delimiter string, values apps.RowValues) string {
	returnStr := ""
	for i, val := range values {
		switch val.(type) {
		case string:
			returnStr += "'"
			returnStr += val.(string)
			returnStr += "'"
		case int:
			returnStr += fmt.Sprintf("%d", val.(int))
		case int64:
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

func PrefixAndSuffixDelimitedSliceStrBuilderSQLRows(prefix string, entries apps.RowEntries, suffix string) string {
	sb := strings.Builder{}
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}

	for count, row := range entries.Rows {

		sb.WriteString("(")
		sb.WriteString(StringDelimitedSliceBuilderSQL(",", row))
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

func DelimitedSliceStrBuilderSQLRows(prefix string, entries apps.RowEntries) string {
	sb := strings.Builder{}
	if len(prefix) > 0 {
		sb.WriteString(prefix)
	}

	for count, row := range entries.Rows {

		sb.WriteString("(")

		for i, val := range row {
			switch val.(type) {
			case string:
				sb.WriteString("'")
				sb.WriteString(val.(string))
				sb.WriteString("'")
			case int, int64:
				returnStr := fmt.Sprintf("%d", val.(int64))
				sb.WriteString(returnStr)
			case uint64:
				returnStr := fmt.Sprintf("%d", val.(uint64))
				sb.WriteString(returnStr)
			case bool:
				returnStr := fmt.Sprintf("%t", val.(bool))
				sb.WriteString(returnStr)
			default:
			}
			if len(row)-1 != i {
				sb.WriteString(",")
			}
		}

		sb.WriteString(")")
		if len(entries.Rows)-1 == count {
			return sb.String()
		}
		sb.WriteString(",")
	}
	return sb.String()
}

func AnyArraySliceStrBuilderSQL(entries apps.RowValues) string {
	var sb strings.Builder

	sb.WriteString("ANY(ARRAY[")

	for i, val := range entries {
		switch val.(type) {
		case string:
			sb.WriteString("'")
			sb.WriteString(val.(string))
			sb.WriteString("'")
		case int:
			returnStr := fmt.Sprintf("%d", val.(int))
			sb.WriteString(returnStr)
		case int64:
			returnStr := fmt.Sprintf("%d", val.(int64))
			sb.WriteString(returnStr)
		case uint64:
			returnStr := fmt.Sprintf("%d", val.(uint64))
			sb.WriteString(returnStr)
		case bool:
			returnStr := fmt.Sprintf("%t", val.(bool))
			sb.WriteString(returnStr)
		default:
		}

		if len(entries)-1 == i {
			sb.WriteString("])")
			return sb.String()
		}
		sb.WriteString(",")
	}

	return sb.String()
}

func ArraySliceStrBuilderSQL(entries apps.RowValues) string {
	var sb strings.Builder

	sb.WriteString("ARRAY[")
	for i, val := range entries {

		switch val.(type) {
		case string:
			sb.WriteString("'")
			sb.WriteString(val.(string))
			sb.WriteString("'")
		case int, int64:
			returnStr := fmt.Sprintf("%d", val.(int64))
			sb.WriteString(returnStr)

		case uint64:
			returnStr := fmt.Sprintf("%d", val.(uint64))
			sb.WriteString(returnStr)

		case bool:
			returnStr := fmt.Sprintf("%t", val.(bool))
			sb.WriteString(returnStr)

		default:
		}

		if len(entries)-1 == i {
			sb.WriteString("]")
			return sb.String()
		}
		sb.WriteString(",")
	}
	sb.WriteString("]")
	return sb.String()
}

func MultiArraySliceStrBuilderSQL(r apps.RowEntries) string {
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
