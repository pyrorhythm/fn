package errs

import "fmt"

func Wrap(err error, s string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", s, err)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

func W[T any](v T, err error) func(string) (T, error) {
	return func(s string) (T, error) {
		return v, Wrap(err, s)
	}
}
