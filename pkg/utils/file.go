package utils

import (
	"bytes"
	"path"
)

type ByteFile struct {
	*bytes.Reader
}

func (b ByteFile) Close() error {
	return nil
}

func NewByteFile(data []byte) ByteFile {
	return ByteFile{Reader: bytes.NewReader(data)}
}

func FormatMediaPathURL(imagePath, publicURL string) string {
	return path.Join(publicURL, imagePath)
}
