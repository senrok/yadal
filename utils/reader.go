package utils

import (
	"io"
	"os"
)

type fileLimit struct {
	reader io.Reader
	close  func() error
}

func (f fileLimit) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f fileLimit) Close() error {
	return f.close()
}

func NewFileLimitReader(file *os.File, size int64) io.ReadCloser {
	return &fileLimit{close: file.Close, reader: io.LimitReader(file, size)}
}
