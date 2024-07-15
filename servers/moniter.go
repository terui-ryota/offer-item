package servers

import (
	"net/http"

	ocprometheus "contrib.go.opencensus.io/exporter/prometheus"
)

type monitorServerOption struct {
	ServerOption
	healthReporters                        []interface{}
	ignoreMetricReportersRegistrationError bool
}

type MonitorOption func(o *monitorServerOption)

type monitorServer struct {
	opt monitorServerOption

	server *http.Server
}

type MetaDataResponse struct {
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	LibGoVersion string `json:"libgo_version"`
}

var ocMetricsExporter *ocprometheus.Exporter

//func (s *monitorServer) Start(ctx context.Context) error {
//	logger.Default().Info("Starting monitor server")
//
//	mux := http.NewServeMux()
//	if ocMetricsExporter != nil {
//		mux.Handle("/metrics", ocMetricsExporter)
//	}
//	mux.Handle(
//		"/live",
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			fmt.Fprint(w, `true`)
//		}))
//	mux.Handle(
//		"/health",
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			var (
//				statusCode = http.StatusOK
//				detail     = make(map[string]interface{})
//			)
//
//			for _, reporter := range s.opt.healthReporters {
//				switch ea := reporter.(type) {
//				case metrics.MapErrorAccessor:
//					stat := make(map[string]bool)
//
//					for k, e := range ea.LastErrors() {
//						isNormal := e == nil
//						stat[k] = isNormal
//						if !isNormal {
//							statusCode = http.StatusInternalServerError
//						}
//					}
//					detail[ea.Name()] = stat
//
//				case metrics.ErrorAccessor:
//					if ea.LastError() != nil {
//						statusCode = http.StatusInternalServerError
//					}
//					detail[ea.Name()] = ea.LastError() == nil
//
//				default:
//					logger.Default().Warnf("skip unknown health reporter %T", reporter)
//				}
//			}
//
//			bs, err := json.Marshal(detail)
//			if err != nil {
//				logger.FromContext(r.Context()).Error("failed to marshal status", zap.Error(err))
//				statusCode = http.StatusInternalServerError
//				bs, _ = json.Marshal(err.Error())
//			}
//
//			w.WriteHeader(statusCode)
//			fmt.Fprint(w, string(bs))
//		}))
//	mux.HandleFunc("/debug/pprof/", pprof.Index)
//	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
//	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
//	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
//	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
//	mux.HandleFunc(
//		"/meta",
//		func(w http.ResponseWriter, r *http.Request) {
//			response := &MetaDataResponse{
//				consts.Version,
//				consts.VcsRevision,
//				consts.LibGoVersion,
//			}
//			res, err := json.Marshal(response)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusInternalServerError)
//				return
//			}
//			w.Header().Set("Content-Type", "application/json")
//			fmt.Fprint(w, string(res))
//		})
//	mux.HandleFunc(
//		"/meta/version",
//		func(w http.ResponseWriter, r *http.Request) {
//			fmt.Fprint(w, consts.Version)
//		})
//	mux.HandleFunc(
//		"/meta/revision",
//		func(w http.ResponseWriter, r *http.Request) {
//			fmt.Fprint(w, consts.VcsRevision)
//		})
//	mux.HandleFunc(
//		"/meta/libgo_version",
//		func(w http.ResponseWriter, r *http.Request) {
//			fmt.Fprint(w, consts.LibGoVersion)
//		},
//	)
//
//	s.server = &http.Server{}
//	s.server.Handler = mux
//
//	for _, reporter := range s.opt.metricReporters {
//		if s.opt.ignoreMetricReportersRegistrationError {
//			for _, c := range reporter.Collector() {
//				if err := prometheus.DefaultRegisterer.Register(c); err != nil {
//					logger.Default().Warn("Ignore MetricsReporter registration error", zap.Error(err))
//				}
//			}
//		} else {
//			prometheus.DefaultRegisterer.MustRegister(reporter.Collector()...)
//		}
//	}
//
//	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opt.port))
//	if err != nil {
//		return fmt.Errorf("failed to listen at %d : %w", s.opt.port, err)
//	}
//
//	return ServeAndWait(
//		"monitor",
//		100*time.Millisecond,
//		5*time.Second,
//		func() error {
//			if err := s.server.Serve(lis); err != nil {
//				return fmt.Errorf("failed to serve monitor server : %w", err)
//			}
//			return nil
//		},
//		testHTTP(s.opt.testHost, s.opt.port, 100*time.Millisecond, nil))
//}
//
//func (s *monitorServer) Stop() error {
//	logger.Default().Info("Stopping monitor server")
//
//	if s.server != nil {
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		if err := s.server.Shutdown(ctx); err != nil {
//			return fmt.Errorf("failed to shutdown monitor server : %w", err)
//		}
//	}
//	for _, reporter := range s.opt.metricReporters {
//		for _, collector := range reporter.Collector() {
//			prometheus.DefaultRegisterer.Unregister(collector)
//		}
//	}
//	return nil
//}
//
//func WithMetricReporter(metricReporters ...metrics.MetricReporter) MonitorOption {
//	return func(o *monitorServerOption) {
//		logger.Default().Debug("Add MetricsReporter")
//		o.metricReporters = append(o.metricReporters, metricReporters...)
//	}
//}
//
//func WithMetricReporterRegistrationErrorIgnoring(ignore bool) MonitorOption {
//	return func(o *monitorServerOption) {
//		logger.Default().Debug("Add WithMetricReporterRegistrationErrorIgnoring", zap.Bool("ignoring", ignore))
//		o.ignoreMetricReportersRegistrationError = ignore
//	}
//}
//
//func WithHealthReporter(healthReporters ...interface{}) MonitorOption {
//	return func(o *monitorServerOption) {
//		logger.Default().Debug("Add MetricsReporter")
//		o.healthReporters = append(o.healthReporters, healthReporters...)
//	}
//}
//
//func NewMonitor(port uint16, opts ...MonitorOption) Server {
//	o := monitorServerOption{
//		ServerOption: ServerOption{
//			testHost: "127.0.0.1",
//			port:     port,
//		},
//		metricReporters: []metrics.MetricReporter{
//			metrics.NewLoggerMetricReporter(),
//		},
//	}
//
//	for _, opt := range opts {
//		opt(&o)
//	}
//
//	return NewMonitorFromOption(o)
//}
//
//func NewMonitorFromOption(opt monitorServerOption) Server {
//	return &monitorServer{opt: opt}
//}
