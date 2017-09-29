package ar

import (
	"strconv"
)

var (
	WHITESPACE_BYTE = []byte(" ")[0]
)

func bytesTrimRight(b []byte) []byte {
	i := len(b) - 1
	for i > 0 && b[i] == WHITESPACE_BYTE {
		i--
	}
	return b[0:i+1]
}

func bytesToString(b []byte) string {
	return string(bytesTrimRight(b))
}

func bytesToDecimal(b []byte) int64 {
	n, _ := strconv.ParseInt(string(bytesTrimRight(b)), 10, 64)
	return n
}

func bytesToOctal(b []byte) int64 {
	n, _ := strconv.ParseInt(string(bytesTrimRight(b)[3:]), 8, 64)
	return n
}
