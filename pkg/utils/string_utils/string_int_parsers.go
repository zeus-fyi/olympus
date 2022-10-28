package string_utils

import (
	"fmt"
	"strconv"
)

func Int64StringParser(str64 string) (int64Value int64) {
	int64Value, err := strconv.ParseInt(str64, 0, 64)

	if err != nil {
		panic(err)
	}
	return int64Value
}

func Uint64StringParser(str64u string) (uint64Value uint64) {
	uint64Value, err := strconv.ParseUint(str64u, 0, 64)

	if err != nil {
		panic(err)
	}
	return uint64Value
}

func IntStringParser(strInt string) (intValue int) {
	intValue, err := strconv.Atoi(strInt)

	if err != nil {
		panic(err)
	}
	return intValue
}

func Convert32BitPtrIntToString(int32BitPtr *int32) string {
	rc := "0"
	if int32BitPtr != nil {
		return rc
	}
	return fmt.Sprintf("%d", int32BitPtr)
}

func ConvertStringTo32BitPtrInt(int32BitPtrString string) *int32 {
	rc := ConvertStringTo32BitInt(int32BitPtrString)
	return &rc
}

func ConvertStringTo32BitInt(int32BitPtrString string) int32 {
	var rc int32
	int32Value, err := strconv.ParseInt(int32BitPtrString, 0, 32)
	if err != nil {
		panic(err)
	}
	rc = int32(int32Value)
	return rc
}
