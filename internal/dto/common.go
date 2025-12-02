package dto

import (
	"io"
)

type FileData struct {
	Name     string
	Mimetype string
	File     io.Reader
}
