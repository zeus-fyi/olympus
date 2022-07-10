package strings

import "strconv"

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
