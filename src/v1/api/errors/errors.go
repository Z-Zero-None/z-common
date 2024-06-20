package errors

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"runtime"
)

const Depth = 32

func New(id int32, code codes.Code, msg string) *Error {
	err := &Error{
		ID:         id,
		Code:       int32(code),
		Message:    msg,
		Stacktrace: []string{},
		Metadata:   map[string]string{},
	}

	var pcs [Depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		f, more := frames.Next()
		err.Stacktrace = append(err.Stacktrace, fmt.Sprintf("%s:%d", f.Function, f.Line))
		if !more {
			break
		}
	}
	//for _, s := range pcs[0:n] {
	//	f := runtime.FuncForPC(s)
	//	_, line := f.FileLine(s - 1)
	//	err.Stacktrace = append(err.Stacktrace, fmt.Sprintf("%s:%d", f.Name(), line))
	//}
	return err
}

func NewFromError(id int32, err error) error {
	if err == nil {
		return nil
	}

	e := FromError(err)
	e.ID = id
	return e
}

func BadRequest(id int32, msg string) error {
	return New(id, http.StatusBadRequest, msg)
}

func Is(err error, targetErr error) bool {
	return errors.Is(err, targetErr)
}

func IsBadRequest(err error) bool {
	return status.Code(err) == http.StatusBadRequest
}

func Unauthorized(id int32, msg string) error {
	return New(id, http.StatusUnauthorized, msg)
}

func IsUnauthorized(err error) bool {
	return status.Code(err) == http.StatusUnauthorized
}

func Forbidden(id int32, msg string) error {
	return New(id, http.StatusForbidden, msg)
}

func IsForbidden(err error) bool {
	return status.Code(err) == http.StatusForbidden
}

func NotFound(id int32, msg string) error {
	return New(id, http.StatusNotFound, msg)
}

func IsNotFound(err error) bool {
	return status.Code(err) == http.StatusNotFound
}

func Conflict(id int32, msg string) error {
	return New(id, http.StatusConflict, msg)
}

func IsConflict(err error) bool {
	return status.Code(err) == http.StatusConflict
}

func ResourceExhausted(id int32, msg string) error {
	return New(id, http.StatusTooManyRequests, msg)
}

func IsResourceExhausted(err error) bool {
	return status.Code(err) == http.StatusTooManyRequests
}

func Internal(id int32, msg string) error {
	return New(id, http.StatusInternalServerError, msg)
}

func IsInternal(err error) bool {
	return status.Code(err) == http.StatusInternalServerError
}

func (e *Error) Error() string {
	return fmt.Sprintf("id=%d, code=%d, message=%s, metadata=%v,stacktrace=%#v", e.ID, e.Code, e.Message, e.Metadata, e.Stacktrace)
}

func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(e.Code), e.Message).WithDetails(e)
	return s
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if ne := new(Error); errors.As(err, &ne) {
		return ne
	}

	st, ok := status.FromError(err)
	ne := New(0, http.StatusInternalServerError, err.Error())
	if !ok {
		return ne
	}
	ne.Code = int32(st.Code())
	ne.Message = st.Message()
	for _, detail := range st.Details() {
		switch e := detail.(type) {
		case *Error:
			ne.ID = e.ID
			ne.Message = e.Message
			ne.Metadata = e.Metadata
			ne.Stacktrace = append(e.Stacktrace[:], ne.Stacktrace[:]...)
		}
	}

	return ne
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	e := FromError(err)

	msg := fmt.Sprintf(format, args...)
	if len(e.Message) > 0 {
		msg = fmt.Sprintf("%s -> %s", msg, e.Message)
	}

	return &Error{
		ID:         e.ID,
		Code:       e.Code,
		Message:    msg,
		Stacktrace: e.Stacktrace,
		Metadata:   e.Metadata,
	}
}
