package acceptance

import (
	"fmt"
	"strconv"
)

// Values holds the values of one at_data.toml block.
type Values map[string]string

func (values Values) String(name string) string {
	value, exists := values[name]
	if !exists {
		panic(fmt.Sprintf("at_data.toml has no value %q", name))
	}
	return value
}

func (values Values) Bool(name string) bool {
	value, err := strconv.ParseBool(values.String(name))
	if err != nil {
		panic(fmt.Sprintf("value %q in at_data.toml: %v", name, err))
	}
	return value
}

// Expect returns an error when actual differs from expected.
func Expect[T comparable](name string, expected, actual T) error {
	if expected != actual {
		return fmt.Errorf("%s: expected %v, got %v", name, expected, actual)
	}
	return nil
}
