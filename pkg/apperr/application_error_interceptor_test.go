package apperr

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	error01 = newAppErr("error01", "error01", codes.NotFound)
	error02 = errors.New("OtherError")
)

func unaryInvoker() grpc.UnaryInvoker {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return error01
	}
}

func TestApplicationErrorUnaryClientInterceptor(t *testing.T) {
	err := ApplicationErrorUnaryClientInterceptor()(context.Background(), "", nil, nil, nil, unaryInvoker(), nil)

	if !xerrors.Is(err, error01) {
		t.Error("error01 should be equal to err")
	}
}

func unaryHandler() grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, error01
	}
}

func TestApplicationErrorUnaryServerInterceptor(t *testing.T) {
	tests := []struct {
		name        string
		occurredErr error
		want        ApplicationError
	}{
		{
			name:        "ApplicationError発生時にエラーが伝播されること",
			occurredErr: newAppErr("error01", "error01", codes.NotFound),
			want:        newAppErr("error01", "error01", codes.NotFound),
		},
		{
			name:        "ApplicationErrorがwrapされているときにエラーが伝播されること",
			occurredErr: fmt.Errorf("error: %w", newAppErr("error01", "error01", codes.NotFound)),
			want:        newAppErr("error01", "error01", codes.NotFound),
		},
		{
			name:        "不明なエラー発生時はUnknownErrorが返却されること",
			occurredErr: errors.New("unexpected error"),
			want:        UnknownError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, tt.occurredErr
			}
			_, err := ApplicationErrorUnaryServerInterceptor()(context.Background(), nil, nil, h)
			var actualError interface{ GRPCStatus() *status.Status }
			if errors.As(err, &actualError) {
				assert.Equal(t, tt.want.GRPCStatus().Code(), actualError.GRPCStatus().Code(), "gRPC status code should be equal")
			} else {
				t.Error("error should implement GRPCStatus()")
			}

			if errors.Is(err, tt.want) {
				t.Errorf("%s, should be equal to %s", err, tt.want)
			}
		})
	}
}

func unaryHandler_OtherError() grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, error02
	}
}

func TestApplicationErrorUnaryServerInterceptor_OtherError(t *testing.T) {
	_, err := ApplicationErrorUnaryServerInterceptor()(context.Background(), nil, nil, unaryHandler_OtherError())

	if e, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
		if e.GRPCStatus().Code() != UnknownError.GRPCStatus().Code() {
			t.Error("error code should be identical with StatusDetailsBuildError")
		}
	} else {
		t.Error("error should implement GRPCStatus()")
	}

	if xerrors.Is(err, UnknownError) {
		t.Error("error01 should be equal to StatusDetailsBuildError")
	}
}

func streamHandler(arg error) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return arg
	}
}

type mockServerStream struct {
	grpc.ServerStream
}

func (s *mockServerStream) Context() context.Context {
	return context.Background()
}

func serverStream() grpc.ServerStream {
	return &mockServerStream{}
}

func TestApplicationErrorStreamServerInterceptor(t *testing.T) {
	tests := []struct {
		name        string
		occurredErr error
		want        ApplicationError
	}{
		{
			name:        "ApplicationError発生時にエラーが伝播されること",
			occurredErr: error01,
			want:        error01,
		},
		{
			name:        "ApplicationErrorがwrapされているときにエラーが伝播されること",
			occurredErr: fmt.Errorf("error: %w", newAppErr("error01", "error01", codes.NotFound)),
			want:        newAppErr("error01", "error01", codes.NotFound),
		},
		{
			name:        "不明なエラー発生時はUnknownErrorが返却されること",
			occurredErr: error02,
			want:        UnknownError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplicationErrorStreamServerInterceptor()(nil, serverStream(), nil, streamHandler(tt.occurredErr))
			var e interface{ GRPCStatus() *status.Status }
			if errors.As(err, &e) {
				if e.GRPCStatus().Code() != tt.want.GRPCStatus().Code() {
					t.Error("error code should be identical")
				}
			} else {
				t.Error("error should implement GRPCStatus()")
			}

			if errors.Is(err, tt.want) {
				t.Errorf("%s, should be equal to %s", err, tt.want)
			}
		})
	}
}

func streamer(arg error) grpc.Streamer {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, arg
	}
}

func TestApplicationErrorStreamClientInterceptor(t *testing.T) {
	tests := []struct {
		name        string
		occurredErr error
		want        ApplicationError
	}{
		{
			name:        "ApplicationError発生時にエラーが伝播されること",
			occurredErr: error01,
			want:        error01,
		},
		{
			name:        "不明なエラー発生時はUnknownErrorが返却されること",
			occurredErr: error02,
			want:        UnknownError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ApplicationErrorStreamClientInterceptor()(context.Background(), nil, nil, "", streamer(tt.occurredErr), nil)

			if !errors.Is(err, tt.want) {
				t.Errorf("%s, should be equal to %s", err, tt.want)
			}
		})
	}
}

func Test_isEnabledErrorLog(t *testing.T) {
	type args struct {
		err                   error
		info                  *grpc.UnaryServerInfo
		ignoreErrorLogMethods map[string][]ApplicationError
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "info が未設定の場合は true",
			args: args{
				err:                   errors.New("err"),
				info:                  nil,
				ignoreErrorLogMethods: map[string][]ApplicationError{},
			},
			want: true,
		},
		{
			name: "ignoreErrorLogMethods が空の場合は true",
			args: args{
				err: errors.New("err"),
				info: &grpc.UnaryServerInfo{
					Server:     serverStream,
					FullMethod: "hoge",
				},
				ignoreErrorLogMethods: map[string][]ApplicationError{},
			},
			want: true,
		},
		{
			name: "ignoreErrorLogMethods の FullMethod は一致するがエラーが異なる場合は true",
			args: args{
				err: AffiliatorNotFoundError,
				info: &grpc.UnaryServerInfo{
					Server:     serverStream,
					FullMethod: "hoge",
				},
				ignoreErrorLogMethods: map[string][]ApplicationError{
					"hoge": {
						AffiliateItemNotFound,
					},
				},
			},
			want: true,
		},
		{
			name: "ignoreErrorLogMethods のエラーは一致するが FullMethod が異なる場合は true",
			args: args{
				err: AffiliatorNotFoundError,
				info: &grpc.UnaryServerInfo{
					Server:     serverStream,
					FullMethod: "fuga",
				},
				ignoreErrorLogMethods: map[string][]ApplicationError{
					"hoge": {
						AffiliatorNotFoundError,
					},
				},
			},
			want: true,
		},
		{
			name: "ignoreErrorLogMethods の FullMethod とエラーが一致する場合は false",
			args: args{
				err: AffiliatorNotFoundError,
				info: &grpc.UnaryServerInfo{
					Server:     serverStream,
					FullMethod: "hoge",
				},
				ignoreErrorLogMethods: map[string][]ApplicationError{
					"hoge": {
						AffiliatorNotFoundError,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEnabledErrorLog(tt.args.err, tt.args.info, tt.args.ignoreErrorLogMethods); got != tt.want {
				t.Errorf("isEnabledErrorLog() = %v, want %v", got, tt.want)
			}
		})
	}
}
