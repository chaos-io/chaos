package errorx

import "fmt"

type Message struct {
	message string
	cause   error
}

func (m *Message) Unwrap() error {
	return m.cause
}

func (m *Message) Error() string {
	return fmt.Sprintf("%s\ncause=%s", m.message, m.cause.Error())
}

func wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return &Message{
		message: fmt.Sprintf(format, args...),
		cause:   err,
	}
}

func WrapfWithStack(err error, format string, args ...any) error {
	return WithStack(wrapf(err, format, args...))
}
