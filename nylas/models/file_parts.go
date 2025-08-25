package models

import "io"

// FilePart describes a file you want to send via multipart with a message.
type FilePart struct {
	Filename    string
	ContentType string // optional (server will default if empty)
	Reader      io.Reader
}
