package uploader

import "io"

type Config struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type Request struct {
	Config

	Image     io.ReadSeekCloser
	Timestamp string
	Size      int64
}

type Uploader interface{}
