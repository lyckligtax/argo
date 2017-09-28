package ar

import (
	"strconv"
	"strings"
)

func stringToBytes(s string, l int) []byte {
	return []byte(growString(s, l))
}

func decimalToBytes(i int64, l int) []byte {
	return []byte(growString(strconv.FormatInt(i, 10), l))
}

func octalToBytes(i int64, l int) []byte {
	return []byte(growString(strconv.FormatInt(i, 8), l))
}

func growString(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}
