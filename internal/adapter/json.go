package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const maxHookInputBytes = 1024 * 1024

func decodeJSON(input io.Reader, target any) error {
	data, err := io.ReadAll(io.LimitReader(input, maxHookInputBytes+1))
	if err != nil {
		return fmt.Errorf("read hook input: %w", err)
	}
	if len(data) > maxHookInputBytes {
		return fmt.Errorf("hook input exceeds %d bytes", maxHookInputBytes)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode hook input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("multiple JSON values are not supported")
		}
		return fmt.Errorf("decode hook input: %w", err)
	}
	return nil
}

func encodeJSON(output io.Writer, value any) error {
	if err := json.NewEncoder(output).Encode(value); err != nil {
		return fmt.Errorf("encode hook output: %w", err)
	}
	return nil
}
