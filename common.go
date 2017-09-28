package ar

import (
	"time"
)

const (
	MAGIC_STRING = "!<arch>\n"
	TRAILER      = "`\n"
	// Lengths in Byte
	HEADER_LENGTH         = 60
	HEADER_OBJECT_NAME    = 16
	HEADER_MTIME          = 12
	HEADER_UID            = 6
	HEADER_GID            = 6
	HEADER_FILE_MODE      = 8
	HEADER_CONTENT_LENGTH = 10
	HEADER_TRAILER        = 2
)

type Header struct {
	ObjectName    string
	MTime         time.Time
	UID           int64
	GID           int64
	FileMode      int64 // Octal
	ContentLength int64
}
