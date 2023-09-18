// Package compressor предоставляет функции для сжатия данных с использованием алгоритма сжатия gzip.
package compressor

import (
	"bytes"
	"compress/gzip"
)

// CompressData сжимает данные перед отправкой серверу с использованием алгоритма сжатия gzip и возвращает сжатые данные в виде буфера.
// Функция принимает на вход срез байтов 'data' и возвращает указатель на bytes.Buffer, содержащий сжатые данные, а также ошибку, если таковая возникла.
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
