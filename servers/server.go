package servers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/terui-ryota/offer-item/pkg/logger"
)

const (
	GrpcLivenessCheckPath = "/liveness.LivenessService/IsLive"
)

const (
	AccessLogFieldRemoteIp   = "remote_ip"
	AccessLogFieldClientId   = "client_id"
	AccessLogFieldUserAgent  = "user_agent"
	AccessLogFieldMethod     = "method" // HTTPOnly
	AccessLogFieldPath       = "path"
	AccessLogFieldStatus     = "status" // StatusCode(int)
	AccessLogFieldStatusCode = "status_code"
	AccessLogFieldDuration   = "duration"
	AccessLogFieldErrorMsg   = "error_msg"
)

// testInterval = 500 * time.Millisecond
// testTimeout  = 5 * time.Second
var ErrorTimeout = errors.New("serveAndWait is timeout")

type Server interface {
	Start(ctx context.Context) error
	Stop() error
}

type ServerOption struct {
	port                     uint16
	testHost                 string
	enableAccessLogger       bool
	accessLoggerExcludes     []string
	disableOpentracing       bool
	healthCheckEndpointPaths map[string]struct{}
	listener                 net.Listener
	startTimeout             time.Duration
}

type Option func(o *ServerOption)

func WithAccessLogger(excludes []string) Option {
	return func(o *ServerOption) {
		logger.Default().Debug("Enable access logger")
		o.enableAccessLogger = true
		o.accessLoggerExcludes = excludes
	}
}

func WithoutOpentracing() Option {
	return func(o *ServerOption) {
		logger.Default().Debug("Disable opentracing")
		o.disableOpentracing = true
	}
}

// WithoutOpentracingForHealthCheckEndpoint は指定された health check endpoint に対して opentracing を無効化します。
// paths: health check endpoint のパス
func WithoutOpentracingForHealthCheckEndpoint(paths []string) Option {
	return func(o *ServerOption) {
		logger.Default().Debug("Disable opentracing for health check endpoint")
		if o.healthCheckEndpointPaths == nil {
			o.healthCheckEndpointPaths = make(map[string]struct{}, len(paths))
		}
		for _, path := range paths {
			o.healthCheckEndpointPaths[path] = struct{}{}
		}
	}
}

func WithTestHost(host string) Option {
	return func(o *ServerOption) {
		o.testHost = host
	}
}

func WithListener(l net.Listener) Option {
	return func(o *ServerOption) {
		o.listener = l
	}
}

func WithStartTimeout(timeout time.Duration) Option {
	return func(o *ServerOption) {
		o.startTimeout = timeout
	}
}

type (
	contextKeyExtraAccessLogFieldType int
	extraAccessLogFieldHolder         struct {
		fields []zap.Field
	}
)

var contextKeyExtraAccessLogField = contextKeyExtraAccessLogFieldType(1)

func newAccessLogFieldHolder(ctx context.Context) (context.Context, *[]zap.Field) {
	holder := &extraAccessLogFieldHolder{fields: []zap.Field{}}
	return context.WithValue(ctx, contextKeyExtraAccessLogField, holder), &holder.fields
}

func AppendAccessLogField(ctx context.Context, fields ...zap.Field) {
	v := ctx.Value(contextKeyExtraAccessLogField)
	if v == nil {
		logger.FromContext(ctx).Error("failed to lookup extra access log holder")
		return
	}
	holder, ok := v.(*extraAccessLogFieldHolder)
	if !ok {
		logger.FromContext(ctx).Error("invalid type for contextKeyExtraAccessLogField")
		return
	}
	holder.fields = append(holder.fields, fields...)
}

// ServeAndWait は serveFn を実行し、 testFn が trueを返すまで 最大 timeout 間 interval 間隔でチェックを行います。
func ServeAndWait(name string, interval, timeout time.Duration, serveFn func() error, testFn func() bool) error {
	var (
		done = errors.New("done")
		snf  = zap.String("server_name", name)
	)

	ch := make(chan error, 1)

	go func() {
		if err := serveFn(); err != nil {
			logger.Default().Info("Starting server", snf)

			ch <- fmt.Errorf("failed to serve : %w", err)
		}
	}()

	go func() {
		logger.Default().Info("Waiting server become ready", snf)

		for {
			time.Sleep(interval)
			logger.Default().Info("Testing server")
			if testFn() {
				logger.Default().Info("Server is now ready", snf)
				ch <- done
				return
			}
		}
	}()

	select {
	case err := <-ch:
		if errors.Is(err, done) {
			return nil
		}
		logger.Default().Error("failed to serve", snf, zap.Error(err))
		return fmt.Errorf("failed to serve : %w", err)
	case <-time.After(timeout):
		logger.Default().Error("timeout")
		return ErrorTimeout
	}
}

func testHTTP(host string, port uint16, interval time.Duration, dialContext func(ctx context.Context, network, addr string) (net.Conn, error)) func() bool {
	return func() bool {
		client := http.Client{Timeout: interval / 2}
		if dialContext != nil {
			client.Transport = &http.Transport{DialContext: dialContext}
		}
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("http://%s:%d", host, port), nil)
		body, err := client.Do(req)
		if err != nil {
			logger.Default().Warn(
				"failed to request to http server", zap.Uint16("port", port), zap.Error(err))
			return false
		}
		if err := body.Body.Close(); err != nil {
			logger.Default().Warn("failed to close response body", zap.Error(err))
		}
		return true
	}
}
