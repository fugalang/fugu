package helper

import (
	"strconv"
)

func BoolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func StringIn64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(2)
	}
	return i64
}

func StringInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func StringFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func StringComplex128(s string) complex128 {
	v, err := strconv.ParseComplex(s, 128)
	if err != nil {
		panic(err)
	}
	return v
}
