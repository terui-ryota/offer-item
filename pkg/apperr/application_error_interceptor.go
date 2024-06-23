package apperr

import (
	"context"
	"errors"

	"github.com/terui-ryota/offer-item/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ApplicationErrorUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return ApplicationErrorUnaryServerInterceptorWithIgnoreErrorLog(map[string][]ApplicationError{})
}

func isEnabledErrorLog(err error, info *grpc.UnaryServerInfo, ignoreErrorLogMethods map[string][]ApplicationError) bool {
	if info == nil {
		return true
	}
	if errs, ok := ignoreErrorLogMethods[info.FullMethod]; ok {
		for _, ignoreErr := range errs {
			if errors.Is(ignoreErr, err) {
				return false
			}
		}
	}
	return true
}

func ApplicationErrorUnaryServerInterceptorWithIgnoreErrorLog(ignoreErrorLogMethods map[string][]ApplicationError) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		res, err := handler(ctx, req)
		if err == nil {
			return res, nil
		}
		if isEnabledErrorLog(err, info, ignoreErrorLogMethods) {
			logger.FromContext(ctx).Error("occurred error", zap.Error(err))
		}

		var ae ApplicationError
		if !errors.As(err, &ae) {
			ae = UnknownError.Wrap(err).(ApplicationError)
		}

		return nil, ae.GRPCStatus().Err()
	}
}

func ApplicationErrorStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err == nil {
			return nil
		}
		logger.FromContext(ss.Context()).Error("occurred error", zap.Error(err))

		var ae ApplicationError
		if !errors.As(err, &ae) {
			ae = UnknownError.Wrap(err).(ApplicationError)
		}

		return ae.GRPCStatus().Err()
	}
}

func ApplicationErrorUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}
		st := status.Convert(err)
		if ae, ok := newAppErrFromGRPCStatus(st); ok {
			return ae
		}
		// envoyまたはclient起因のエラーのハンドリング
		if st.Code() == codes.Unavailable {
			return UnavailableError.Wrap(st.Err())
		}
		if st.Code() == codes.Canceled {
			return RequestCanceled.Wrap(st.Err())
		}
		// 不明なエラー
		return UnknownError.Wrap(err)
	}
}

func ApplicationErrorStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err == nil {
			return stream, nil
		}
		st := status.Convert(err)
		if ae, ok := newAppErrFromGRPCStatus(st); ok {
			return stream, ae
		}
		// envoyまたはclient起因のエラーのハンドリング
		if st.Code() == codes.Unavailable {
			return stream, UnavailableError.Wrap(st.Err())
		}
		if st.Code() == codes.Canceled {
			return stream, RequestCanceled.Wrap(st.Err())
		}
		return stream, UnknownError.Wrap(err)
	}
}
