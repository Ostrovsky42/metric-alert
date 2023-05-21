package helpers

import (
	"bytes"
	"encoding/json"
	"io"
)

func EncodeData(data interface{}) (*bytes.Buffer, error) {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func UnmarshalBody(reader io.ReadCloser, str any) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = json.Unmarshal(body, &str)
	if err != nil {
		return err
	}

	return nil
}
