package apperr

import (
	"fmt"

	commonpb "github.com/ca-media-nantes/pick/protofiles/go/common"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApplicationError interface {
	Code() string
	Message() string
	Error() string
	Wrap(err error) error
	Unwrap() error
	GRPCStatus() *status.Status
	Proto() *commonpb.ApplicationError
	Is(err error) bool
}

type appErr struct {
	code    string
	message string
	// ラップするエラー
	err            error
	frame          xerrors.Frame
	gRPCStatusCode codes.Code
}

func newAppErr(code, message string, gRPCStatusCode codes.Code) ApplicationError {
	return &appErr{
		code:           code,
		message:        message,
		err:            nil,
		frame:          xerrors.Caller(1),
		gRPCStatusCode: gRPCStatusCode,
	}
}

func newAppErrFromGRPCStatus(s *status.Status) (ApplicationError, bool) {
	for _, d := range s.Details() {
		switch t := d.(type) {
		case *commonpb.ApplicationError:
			return newAppErr(t.GetCode(), t.GetMessage(), s.Code()), true
		}
	}
	return nil, false
}

// ここでは appErr はポインターにしない
func (e appErr) Wrap(next error) error {
	e.err = next
	e.frame = xerrors.Caller(1)
	return &e
}

func (e *appErr) Code() string {
	return e.code
}

func (e *appErr) Message() string {
	return e.message
}

func (e *appErr) Error() string {
	if e.err != nil {
		return fmt.Sprintf("code: %s, message: %s, caused by:\n%s", e.code, e.message, e.err.Error())
	}

	return fmt.Sprintf("code: %s, message: %s", e.code, e.message)
}

func (e *appErr) GRPCStatus() *status.Status {
	st := status.New(e.gRPCStatusCode, e.Error())
	if st, _ := st.WithDetails(e.Proto()); st != nil {
		return st
	}
	return st
}

func (e *appErr) Proto() *commonpb.ApplicationError {
	return &commonpb.ApplicationError{
		Message: e.message,
		Code:    e.code,
	}
}

func (e *appErr) Unwrap() error {
	return e.err
}

// fmt.Formatterを実装
func (e *appErr) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

// xerrors.Formatterを実装
func (e *appErr) FormatError(p xerrors.Printer) (next error) {
	p.Print(e.Error())
	e.frame.Format(p)
	return e.err
}

func (e *appErr) Is(err error) bool {
	var ae *appErr
	return xerrors.As(err, &ae) && e.code == ae.code
}
