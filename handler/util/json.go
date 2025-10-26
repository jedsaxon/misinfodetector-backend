package util

import (
	"encoding/json"
	"io"
)

func UnmarshalJsonReader(reader io.Reader, payload any) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, &payload); err != nil {
		return err
	}

	return nil
}
