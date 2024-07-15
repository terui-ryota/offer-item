package grpc_proxyserver

import (
	"context"
	"net/http"
	"net/textproto"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opencensus.io/plugin/ochttp"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	//"github.com/ca-media-nantes/pick/go-lib/internal/gateway"
)

const (
	xClientUserAgent = "X-Client-User-Agent"
	xForwardedFor    = "X-Forwarded-For"
)

type grpcProxyServer struct {
	server *http.Server
	conn   *grpc.ClientConn
}

type RegisterProxyHandlerFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

func NewGrpcProxyServer(backendEndpoint string, registerFuncs []RegisterProxyHandlerFunc, checkPermissionFuncs []CheckPermissionFunc, opts ...Option) (GrpcProxyServer, error) {
	opt := defaultOption.clone()
	for _, f := range opts {
		f(opt)
	}
	conn, err := grpc.Dial(
		backendEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithStatsHandler(&ocgrpc.ClientHandler{
		// 	StartOptions: trace.StartOptions{
		// 		Sampler: trace.ProbabilitySampler(0.5),
		// 	},
		// }),
		//grpc.WithUnaryInterceptor(authUnaryServerInterceptor(checkPermissionFuncs...)),
	)
	if err != nil {
		return nil, xerrors.Errorf("failed connect backend server: %w", err)
	}
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
		runtime.WithMetadata(endClientUserAgentAnnotator),
		runtime.WithMetadata(endClientIPAddressAnnotator),
		runtime.WithErrorHandler(HttpErrorHandler),
	)
	for _, f := range registerFuncs {
		if err := f(context.Background(), gwMux, conn); err != nil {
			return nil, xerrors.Errorf("failed register proxy handler: %w", err)
		}
	}

	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(gwMux)
	router.HandleFunc(opt.HealthCheckEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}).Methods("GET")

	server := &http.Server{
		Addr:              opt.Address,
		WriteTimeout:      opt.WriteTimeout,
		ReadTimeout:       opt.ReadTimeout,
		IdleTimeout:       opt.IdleTimeout,
		ReadHeaderTimeout: opt.ReadHeaderTimeout,
		Handler: &ochttp.Handler{
			Handler: router,
		},
	}
	return &grpcProxyServer{
		server: server,
		conn:   conn,
	}, nil
}

func (s *grpcProxyServer) Start() error {
	s.server.RegisterOnShutdown(func() {
		s.server.SetKeepAlivesEnabled(false)
	})
	return s.server.ListenAndServe()
}

func (s *grpcProxyServer) Shutdown(ctx context.Context) error {
	// close grpc connection
	defer s.conn.Close()
	return s.server.Shutdown(ctx)
}

type GrpcProxyServer interface {
	Start() error
	Shutdown(ctx context.Context) error
}

func incomingHeaderMatcher(key string) (string, bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	//if _, ok := gateway.SpecificHeaders[key]; ok {
	//	return key, true
	//}
	return runtime.DefaultHeaderMatcher(key)
}

func endClientUserAgentAnnotator(ctx context.Context, req *http.Request) metadata.MD {
	ua := req.Header.Get(xClientUserAgent)

	return metadata.New(map[string]string{
		xClientUserAgent: ua,
	})
}

func endClientIPAddressAnnotator(ctx context.Context, req *http.Request) metadata.MD {
	ff := req.Header.Get(xForwardedFor)

	return metadata.New(map[string]string{
		xForwardedFor: ff,
	})
}
