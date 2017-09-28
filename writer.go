package ar

import (
	"io"
	"strings"
	"strconv"
	"errors"
	"fmt"
)

var (
	ErrWriteLength    = errors.New("archive/ar: content length does not match")
	ErrWriteAfterClose = errors.New("archive/ar: write after close")
	ErrHeaderTooLong      = errors.New("archive/ar: header too long")
)

type Writer struct {
	w      io.Writer
	closed bool
	nb     int64
}

func NewWriter(w io.Writer) (*Writer, error) {
	arw := &Writer{
		w: w,
	}
	return arw, arw.writeGlobalHeader()
}

func (arw *Writer) Close() {
	arw.closed = true
}

func (arw *Writer) writeGlobalHeader() error {
	_, err := arw.w.Write([]byte(MAGIC_STRING))
	return err
}

func (arw *Writer) Write(content []byte) (n int, err error) {
	if arw.closed {
		return 0, ErrWriteAfterClose
	}

	if arw.nb != int64(len(content)) {
		return 0, ErrWriteLength
	}

	n, err = arw.w.Write(content)
	if err != nil {
		return
	}

	arw.nb -= int64(n)

	if len(content)%2 == 1 {
		l, err := arw.w.Write([]byte{'\n'})
		arw.nb--
		return n + l, err
	}

	return
}

func (arw *Writer) WriteHeader(header Header) (n int, err error) {
	if arw.closed {
		return 0, ErrWriteAfterClose
	}

	if arw.nb != 0 {
		return 0, fmt.Errorf("archive/ar: missed writing %d bytes", arw.nb)
	}

	position := 0
	hdr := make([]byte, HEADER_LENGTH)

	for _, bytes := range [][]byte{
		arString(header.ObjectName, HEADER_OBJECT_NAME),
		arDecimal(header.MTime.Unix(), HEADER_MTIME),
		arDecimal(header.UID, HEADER_UID),
		arDecimal(header.GID, HEADER_GID),
		arOctal(header.FileMode, HEADER_FILE_MODE),
		arDecimal(header.ContentLength, HEADER_CONTENT_LENGTH),
		arString(TRAILER, HEADER_TRAILER),
	} {
		copy(hdr[position:], bytes)
		position += len(bytes)
	}

	arw.nb = header.ContentLength
	n, err = arw.w.Write(hdr)
	if n > HEADER_LENGTH {
		err = ErrHeaderTooLong
	}

	return
}

func arString(s string, l int) []byte {
	return []byte(growString(s, l))
}

func arDecimal(i int64, l int) []byte {
	return []byte(growString(strconv.FormatInt(i, 10), l))
}

func arOctal(i int64, l int) []byte {
	return []byte(growString(strconv.FormatInt(i, 8), l))
}

// growString fills a string s with c until the length l is achieved
func growString(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}
