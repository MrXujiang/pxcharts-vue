package errorx

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type Error struct {
	Status  int
	Message string
}

// Define alias
var (
	WithStack = errors.WithStack
	Wrap      = errors.Wrap
	Wrapf     = errors.Wrapf
	Is        = errors.Is
	Errorf    = errors.Errorf
)

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var reqErr *Error
	if errors.As(err, &reqErr) {
		return reqErr, true
	}
	return nil, false
}

func New(code int, message string) *Error {
	return &Error{
		Status:  code,
		Message: message,
	}
}

func BadRequest(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusBadRequest)
	}
	return New(http.StatusBadRequest, msg)
}

func Unauthorized(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusUnauthorized)
	}
	return New(http.StatusUnauthorized, msg)
}

func Forbidden(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusForbidden)
	}
	return New(http.StatusForbidden, msg)
}

func NotFound(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusNotFound)
	}
	return New(http.StatusNotFound, msg)
}

func MethodNoAllowed(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusMethodNotAllowed)
	}
	return New(http.StatusMethodNotAllowed, msg)
}

func TooManyRequests(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusTooManyRequests)
	}
	return New(http.StatusTooManyRequests, msg)
}

func RequestTimeout(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusRequestTimeout)
	}
	return New(http.StatusRequestTimeout, msg)
}

func Conflict(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusConflict)
	}
	return New(http.StatusConflict, msg)
}

func RequestEntityTooLarge(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusRequestEntityTooLarge)
	}
	return New(http.StatusRequestEntityTooLarge, msg)
}

func InternalServerError(msg string) *Error {
	if msg == "" {
		msg = http.StatusText(http.StatusInternalServerError)
	}
	return New(http.StatusInternalServerError, msg)
}
