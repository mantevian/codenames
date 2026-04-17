package enums

import (
	"database/sql/driver"
	"fmt"
)

type enumBase[T ~string] struct {
	value T
}

func (e *enumBase[T]) Scan(value any) error {
	if value == nil {
		e.value = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		e.value = T(v)
	case []byte:
		e.value = T(string(v))
	default:
		return fmt.Errorf("cannot scan %T into enum", value)
	}
	return nil
}

func (e enumBase[T]) Value() (driver.Value, error) {
	return string(e.value), nil
}
