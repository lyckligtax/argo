package ar

import (
	"errors"
	"io"
	"time"
)

var (
	ErrReadLength     = errors.New("archive/ar: content length does not match")
	ErrReadAfterClose = errors.New("archive/ar: read after close")
)

type Reader struct {
	r      io.Reader
	closed bool
	nb     int64
	pad    int64
	seeker io.Seeker
}

func NewReader(r io.Reader) *Reader {
	if seeker, ok := r.(io.Seeker); ok {
		return &Reader{
			r:      r,
			closed: false,
			nb:     0,
			pad:    int64(len(MAGIC_STRING)),
			seeker: seeker,
		}
	}

	return nil
}

func (arr *Reader) Next() (*Header, error) {
	_, err := arr.seeker.Seek(arr.nb+arr.pad, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return arr.readHeader()
}

// Read data from the current entry in the archive.
func (arr *Reader) Read(content []byte) (n int, err error) {
	if arr.closed {
		return 0, ErrReadAfterClose
	}
	if arr.nb == 0 {
		return 0, io.EOF
	}

	if int64(len(content)) != arr.nb {
		return 0, ErrReadLength
	}
	n, err = arr.r.Read(content)
	arr.nb -= int64(n)

	return
}

func (arr *Reader) readHeader() (*Header, error) {
	walker := NewWalker(make([]byte, HEADER_LENGTH))
	if _, err := io.ReadFull(arr.r, walker.Bytes); err != nil {
		return nil, err
	}

	header := &Header{
		ObjectName:    bytesToString(walker.Next(HEADER_OBJECT_NAME)),
		MTime:         time.Unix(bytesToDecimal(walker.Next(HEADER_MTIME)), 0),
		UID:           bytesToDecimal(walker.Next(HEADER_UID)),
		GID:           bytesToDecimal(walker.Next(HEADER_GID)),
		FileMode:      bytesToOctal(walker.Next(HEADER_FILE_MODE)),
		ContentLength: bytesToDecimal(walker.Next(HEADER_CONTENT_LENGTH)),
	}

	arr.nb = header.ContentLength
	if header.ContentLength%2 == 1 {
		arr.pad = 1
	}

	return header, nil
}
