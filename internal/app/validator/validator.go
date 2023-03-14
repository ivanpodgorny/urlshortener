package validator

import (
	"fmt"
	"net/url"
)

type Validator[T any] func(val T) error

// Validate применяет набор валидаторов validators к значению val.
func Validate[T any](val T, validators ...Validator[T]) (bool, error) {
	for _, v := range validators {
		if err := v(val); err != nil {
			return false, err
		}
	}

	return true, nil
}

// IsURL проверяет, что строка является валидным URL.
func IsURL(val string) error {
	u, err := url.Parse(val)
	if err == nil && u.Scheme != "" && u.Host != "" {
		return nil
	}

	return fmt.Errorf("%s is not valid url", val)
}

// Length возвращает валидатор, который проверяет, что строка не превышает длину l.
func Length(l int) Validator[string] {
	return func(val string) error {
		if len([]rune(val)) > l {
			return fmt.Errorf("%s is too long, %d is max length", val, l)
		}

		return nil
	}
}

// Size возвращает валидатор, который проверяет, что размер массива
// T элементов не превышает s.
func Size[T any](s int) Validator[[]T] {
	return func(val []T) error {
		if len(val) > s {
			return fmt.Errorf("slice has too many elements, %d is max length", s)
		}

		return nil
	}
}
