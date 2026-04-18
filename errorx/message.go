package errorx

import "fmt"

type messageError struct {
	message string
	cause   error
}

func (m *messageError) Unwrap() error {
	return m.cause
}

func (m *messageError) Error() string {
	return fmt.Sprintf("%s\ncause=%s", m.message, m.cause.Error())
}

func wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return &messageError{
		message: fmt.Sprintf(format, args...),
		cause:   err,
	}
}

func WrapfWithStack(err error, format string, args ...any) error {
	return WithStack(wrapf(err, format, args...))
}
