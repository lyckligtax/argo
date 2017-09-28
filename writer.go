package ar

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrWriteLength     = errors.New("archive/ar: content length does not match")
	ErrWriteAfterClose = errors.New("archive/ar: write after close")
	ErrHeaderTooLong   = errors.New("archive/ar: header too long")
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
		stringToBytes(header.ObjectName, HEADER_OBJECT_NAME),
		decimalToBytes(header.MTime.Unix(), HEADER_MTIME),
		decimalToBytes(header.UID, HEADER_UID),
		decimalToBytes(header.GID, HEADER_GID),
		octalToBytes(header.FileMode, HEADER_FILE_MODE),
		decimalToBytes(header.ContentLength, HEADER_CONTENT_LENGTH),
		stringToBytes(TRAILER, HEADER_TRAILER),
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
