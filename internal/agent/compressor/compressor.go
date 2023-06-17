package compressor

import (
	"bytes"
	"compress/gzip"
)

func CompressData(data []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
