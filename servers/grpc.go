package servers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/terui-ryota/offer-item/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
)

type grpcServer struct {
	opt grpcServerOption

	grpcServer *grpc.Server

	cbacClientCnx *grpc.ClientConn
}

type GrpcServiceRegister func(s *grpc.Server)
type grpcServerOption struct {
	ServerOption

	rawServerOptions []grpc.ServerOption

	unaryAuthInterceptor grpc.UnaryServerInterceptor

	serviceRegisters  []GrpcServiceRegister
	unaryInterceptors []grpc.UnaryServerInterceptor

	propagatedHeaders []string

	defaultContextLogFields map[string]string
}

type GrpcOption func(o *grpcServerOption) error

func WithGrpcService(register GrpcServiceRegister) GrpcOption {
	return func(o *grpcServerOption) error {
		logger.Default().Debug("Add GrpcServiceRegister")
		if o.serviceRegisters == nil {
			o.serviceRegisters = make([]GrpcServiceRegister, 0)
		}
		o.serviceRegisters = append(o.serviceRegisters, register)
		return nil
	}
}

func WithGrpcRawServerOption(opts ...grpc.ServerOption) GrpcOption {
	return func(o *grpcServerOption) error {
		o.rawServerOptions = opts
		return nil
	}
}

func WithGrpcUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) GrpcOption {
	return func(o *grpcServerOption) error {
		logger.Default().Debug("Add GrpcUnaryInterceptor")
		if o.unaryInterceptors == nil {
			o.unaryInterceptors = make([]grpc.UnaryServerInterceptor, 0)
		}
		o.unaryInterceptors = append(o.unaryInterceptors, interceptor...)
		return nil
	}
}

func NewGrpcServer(grpc uint16, opts ...interface{}) (Server, error) {
	gopt := grpcServerOption{ServerOption: ServerOption{port: grpc}}
	for _, opt := range opts {
		if copt, ok := opt.(Option); ok {
			copt(&gopt.ServerOption)
		} else if copt, ok := opt.(GrpcOption); ok {
			if err := copt(&gopt); err != nil {
				return nil, fmt.Errorf("failed to configure server : %w", err)
			}
		} else {
			return nil, fmt.Errorf("unhandled option type(%T)", opt)
		}
	}

	return &grpcServer{
		opt: gopt,
	}, nil
}

func (s *grpcServer) Start(ctx context.Context) error {
	var err error
	logger.Default().Info("Starting gRPC server")

	serverOptions := make([]grpc.ServerOption, 0)

	// 生オプションを設定
	serverOptions = append(serverOptions, s.opt.rawServerOptions...)

	// Register interceptor
	interceptors := make([]grpc.UnaryServerInterceptor, 0)

	if s.opt.defaultContextLogFields != nil {
		interceptors = append(interceptors,
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				m := make(map[string]string)
				for k, v := range s.opt.defaultContextLogFields {
					m[k] = v
				}

				ctx = logger.SetupContextFieldsHolder(ctx, m)
				return handler(ctx, req)
			},
		)
	}

	serverOptions =
		append(
			serverOptions,
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(interceptors...)))

	// Generate server
	s.grpcServer =
		grpc.NewServer(serverOptions...)

	// Use zap logger
	grpc_zap.ReplaceGrpcLoggerV2(logger.Default().Native())

	// Register service
	for _, register := range s.opt.serviceRegisters {
		register(s.grpcServer)
	}

	// Register grpc reflection
	reflection.Register(s.grpcServer)

	var lis net.Listener
	if s.opt.listener != nil {
		lis = s.opt.listener
	} else {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", s.opt.port))
		if err != nil {
			return fmt.Errorf("failed to listen at %d : %w", s.opt.port, err)
		}
	}

	startTimeout := 5 * time.Second
	if s.opt.startTimeout > 0 {
		startTimeout = s.opt.startTimeout
	}

	if err :=
		ServeAndWait(
			"gRPC",
			100*time.Millisecond,
			startTimeout,
			func() error {
				if err := s.grpcServer.Serve(lis); err != nil {
					return fmt.Errorf("failed to serve grpc server : %w", err)
				}
				return nil
			},
			func() bool {
				var opts []grpc.DialOption
				opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if lis, ok := s.opt.listener.(*bufconn.Listener); ok {
					opts = append(opts, grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
						return lis.Dial() //nolint:wrapcheck
					}))
				}
				cc, err := grpc.Dial(fmt.Sprintf("%s:%d", s.opt.testHost, s.opt.port), opts...)
				if err != nil {
					logger.Default().Warn("failed to dial to grpc server", zap.Error(err))
					return false
				}
				defer cc.Close()

				return true
			}); err != nil {
		return err
	}

	return nil
}

func (s *grpcServer) Stop() error {
	logger.Default().Info("Stopping gRPC server")

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if s.cbacClientCnx != nil {
		if err := s.cbacClientCnx.Close(); err != nil {
			logger.Default().Error("failed to close cbac client connection", zap.Error(err))
		}
	}

	return nil
}
