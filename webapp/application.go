package webapp

import (
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pomelo-la/go-toolkit/httprouter"
	"github.com/pomelo-la/go-toolkit/logger"
	"github.com/pomelo-la/go-toolkit/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
)

const (
	_defaultWebApplicationPort = "8080"
	_defaultRuntimeEnvironment = "local"

	// Default compression level for defined response content types.
	// The level should be one of the ones defined in the flat package.
	// Higher levels typically run slower but compress more.
	_defaultCompressionLevel = 5
	_requestIDHeader         = "x-request-id"
	_debugHeader             = "x-debug"
)

var _defaultApplicationName = os.Getenv("OTEL_SERVICE_NAME")

// Application is a container struct that contains all required base components
// for building web applications.
type Application struct {
	config AppOptions

	Router  *httprouter.Router
	Runtime Runtime
	Logger  logger.Logger
	Tracer  telemetry.Trace
	Meter   telemetry.Metric
}

// AppOptions represents the options for configuring a web application.
type AppOptions struct {
	ServerTimeouts httprouter.Timeouts
	LogLevel       logger.Level
	Listener       net.Listener
	Runtime        string
	ErrorHandler   httprouter.ErrorHandlerFunc
}

// WithTimeouts allows you to configure the different timeouts
// that the http server uses.
//
// Default behavior is to not have timeouts for incoming requests.
func WithTimeouts(timeout httprouter.Timeouts) func(options *AppOptions) {
	return func(opts *AppOptions) {
		opts.ServerTimeouts = timeout
	}
}

// WithErrorHandler allows you to set a custom error handling function.
//
// The function gets called everytime one of your handlers returns en non-nil error.
// Default is to treat all errors that are not httprouter.Error as 500 status code errors.
func WithErrorHandler(errHandlerFunc httprouter.ErrorHandlerFunc) func(options *AppOptions) {
	return func(opts *AppOptions) {
		opts.ErrorHandler = errHandlerFunc
	}
}

// WithLogLevel allows you to configure the level at which
// the backend logger will log.
//
// Default behavior is to log at Warn level in production, and Debug level in
// local and test environments.
func WithLogLevel(level logger.Level) func(options *AppOptions) {
	return func(opts *AppOptions) {
		opts.LogLevel = level
	}
}

// WithListener allows you to configure the network listener at which the web
// server will be listening to incoming connections.
//
// Default behavior is to use whatever value is in PORT env variable, and if
// none is found, then use 8080.
func WithListener(listener net.Listener) func(options *AppOptions) {
	return func(opts *AppOptions) {
		opts.Listener = listener
	}
}

// WithEnvironmentRuntime allows you to configure the scope string to use for parsing and
// bootstrapping the http server.
//
// Default behavior is to use whatever value is in RUNTIME env variable, and if
// none is found, then assume backend is running locally.
func WithEnvironmentRuntime(environment string) func(options *AppOptions) {
	return func(opts *AppOptions) {
		opts.Runtime = environment
	}
}

// Run starts your Application, it blocks until os.Interrupt is received.
func (a *Application) Run() error {
	err := a.configureListener()
	if err != nil {
		return err
	}

	a.Logger.Info(context.Background(), "http server listening")

	if err := a.printRoutes(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go expvarPolling(ctx)

	// Run blocks until the web backend was signaled to close
	if err := httprouter.Run(a.config.Listener, a.config.ServerTimeouts, a.Router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.Logger.Error(ctx, "startup", "error_msg", err)
		return err
	}

	if err := a.Meter.ShutdownMetricProvider(ctx); err != nil {
		a.Logger.Error(ctx, "shutdown telemetry meter", "error_msg", err)
		return err
	}

	return a.Tracer.ShutdownTraceProvider(ctx)
}

func (a *Application) configureListener() error {
	if a.config.Listener != nil {
		return nil
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = _defaultWebApplicationPort
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	a.config.Listener = ln

	return nil
}

// printRoutes prints every route grouped by URL and http methods.
// Example:
//
// /path                  [GET POST]
// /path/sub-path         [GET]
// /path/{id}             [POST]
// /ping                  [GET].
func (a *Application) printRoutes() error {
	var w tabwriter.Writer
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.TabIndent)

	routes, err := a.Router.Routes()
	if err != nil {
		return err
	}

	m := make(map[string][]string)
	var r []string
	for _, route := range routes {
		r = append(r, route.Route)
		m[route.Route] = append(m[route.Route], route.Method)
	}

	visited := make(map[string]struct{})
	sort.Strings(r)

	for _, v := range r {
		if _, ok := visited[v]; !ok {
			sort.Strings(m[v])
			fmt.Fprintf(&w, "%s\t%v\t\n", v, m[v])
			visited[v] = struct{}{}
		}
	}

	// Flush routes buffer
	return w.Flush()
}

// New instantiates a backend Application with sane defaults.
func New(optFns ...func(opts *AppOptions)) (*Application, error) {
	var config AppOptions
	for _, fn := range optFns {
		fn(&config)
	}

	config.LogLevel = logger.LevelInfo

	runtime, err := configRuntime(config)
	if err != nil {
		return nil, err
	}

	if config.ServerTimeouts == (httprouter.Timeouts{}) {
		config.ServerTimeouts = httprouter.Timeouts{
			ShutdownTimeout: 5 * time.Second,
		}
	}

	log := configureLogger(config)

	if !strings.EqualFold(runtime.Environment, _defaultRuntimeEnvironment) {
		tracer, err := telemetry.NewTrace(context.Background())
		if err != nil {
			return nil, err
		}
		meter, err := telemetry.NewMetric(context.Background())
		if err != nil {
			return nil, err
		}
		router := defaultHTTPRouter(*log, *tracer, config.ErrorHandler)

		return &Application{
			config:  config,
			Router:  router,
			Runtime: *runtime,
			Logger:  *log,
			Tracer:  *tracer,
			Meter:   *meter,
		}, nil
	}

	router := defaultHTTPRouter(*log, telemetry.Trace{}, config.ErrorHandler)
	return &Application{
		config:  config,
		Router:  router,
		Runtime: *runtime,
		Logger:  *log,
		Tracer:  telemetry.Trace{},
		Meter:   telemetry.Metric{},
	}, nil
}

func configureLogger(config AppOptions) *logger.Logger {
	var log *logger.Logger
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "SEND ALERT")
		},
	}
	traceIDFn := func(ctx context.Context) string {
		traceID, err := telemetry.GetTraceID(ctx)
		if err != nil {
			return ""
		}
		return traceID
	}

	return logger.NewWithEvents(os.Stdout, config.LogLevel, _defaultApplicationName, traceIDFn, events)
}

func configRuntime(opt AppOptions) (*Runtime, error) {
	runtime := Runtime{Environment: _defaultRuntimeEnvironment}
	if len(opt.Runtime) == 0 {
		env, err := RuntimeFromEnv()
		if err != nil {
			return nil, err
		}

		runtime = env
	}

	if len(opt.Runtime) != 0 {
		env, err := RuntimeFromString(opt.Runtime)
		if err != nil {
			return nil, err
		}

		runtime = env
	}

	return &runtime, nil
}

func defaultHTTPRouter(log logger.Logger, trace telemetry.Trace, errorHandlerFunc httprouter.ErrorHandlerFunc) *httprouter.Router {
	notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := httprouter.NewErrorf(http.StatusNotFound, "resource %s not found", r.URL.Path)
		_ = httprouter.RespondJSON(w, http.StatusNotFound, err)
	})

	livenessHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = httprouter.RespondJSON(w, http.StatusNoContent, nil)
	})

	readinessHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = httprouter.RespondJSON(w, http.StatusNoContent, nil)
	})

	mdwl := []httprouter.Middleware{
		telemetryMiddleware(trace),
		logMiddleware(log),
		panicsMiddleware(log),
		headerForwarder(trace),
		newCompressor(),
	}

	return httprouter.New(httprouter.Config{
		NotFoundHandler:             notFoundHandler,
		HealthCheckLivenessHandler:  livenessHandler,
		HealthCheckReadinessHandler: readinessHandler,
		EnableProfiler:              true,
		Mw:                          mdwl,
		ErrorHandlerFunc:            errorHandlerFunc,
	})
}

// headerForwarder decorates a request context with the value of certain headers
// in order to allow transport.HTTPRequester to use those headers in outgoing requests.
func headerForwarder(tracer telemetry.Trace) httprouter.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			_, span := tracer.AddSpan(ctx, "HeaderForwarder")
			defer span.End()

			propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
			ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

			span.AddEvent("HeaderForwarder processing")

			r2 := r.WithContext(ctx)
			handler(w, r2)
		}
	}
}

// log decorates the request context with the given logger, accessible via
// the go-core log methods with context.
func logMiddleware(log logger.Logger) httprouter.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(r.Context(), "request started",
				slog.String("method", r.Method),
				slog.String("path", path),
				slog.String("remoteaddr", r.RemoteAddr),
			)

			handler(w, r)
		}
	}
}

// newCompressor returns a middleware that compresses response body of a given content type to a data format based
// on Accept-Encoding request header. It uses the _defaultCompressionLevel.
//
// NOTE: if you don't use web.RespondJSON to marshal the body into the writer,
// make sure to set the Content-Type header on your response otherwise this middleware will not compress the response body.
func newCompressor() httprouter.Middleware {
	c := middleware.NewCompressor(_defaultCompressionLevel)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return c.Handler(next).ServeHTTP
	}
}

// panicsMiddleware handles any panic that may occur by notifying the error to an external telemetry system such NewRelic
// and responding to the client with an `Error` and status code 500.
// For this middleware to log, it requires the context to have a log.Logger.
func panicsMiddleware(log logger.Logger) httprouter.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(r.Context(), "panic recover")

					statusCode := http.StatusInternalServerError
					httprouter.NotifyErr(r, fmt.Errorf("%v", err), statusCode)
					_ = httprouter.RespondJSON(w, statusCode, httprouter.NewErrorf(statusCode, fmt.Sprintf("%v", err)))
				}
			}()

			handler(w, r)
		}
	}
}

// Telemetry middleware simplifies tracing of incoming web requests by
// initiating a new Span and composing the request context with it.
func telemetryMiddleware(tracer telemetry.Trace) httprouter.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			routePattern := chi.RouteContext(r.Context()).RoutePattern()

			txName := fmt.Sprintf("%s (%s)", routePattern, r.Method)

			ctx, span := tracer.AddSpan(r.Context(), txName)
			defer span.End()

			// The Span returned by tracer.StartWebSpan may be a
			// http.ResponseWriter as well. In this case we want to use it when
			// calling the user handler.
			spanWriter, ok := span.(http.ResponseWriter)
			if ok {
				w = spanWriter
			}

			r2 := r.WithContext(ctx)

			// Wrap the http.ResponseWriter with a proxy for later response
			// inspection.
			w2 := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			handler(w2, r2)
			recordRequest(r2.Context(), w2.Status(), time.Since(start), r.Method, routePattern)
		}
	}
}

func recordRequest(ctx context.Context, status int, delta time.Duration, method, routePattern string) {
	meter := otel.Meter("go-toolkit")

	statusClass := strconv.Itoa(status/100) + "xx" // 2xx, 3xx, 4xx, 5xx
	attr := []attribute.KeyValue{
		{
			Key:   "status",
			Value: attribute.IntValue(status),
		},
		{
			Key:   "status_class",
			Value: attribute.StringValue(statusClass),
		},
		{
			Key:   "method",
			Value: attribute.StringValue(method),
		},
		{
			Key:   "handler",
			Value: attribute.StringValue(telemetry.SanitizeMetricTagValue(routePattern)),
		},
	}
	reqIncr, err := meter.Int64UpDownCounter("http.server.request.counter")
	if err != nil {
		return
	}
	reqIncr.Add(ctx, 1, metric.WithAttributes(attr...))

	reqHistogram, err := meter.Float64Histogram("http.server.request.duration")
	if err != nil {
		return
	}

	reqHistogram.Record(ctx, delta.Seconds(), metric.WithAttributes(attr...))
}

type pooledTransportPoolInfo map[string]map[string]int64

func expvarPolling(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			exportedVarPoolHTTP()
		case <-ctx.Done():
			return
		}
	}
}

func exportedVarPoolHTTP() {
	meter := otel.Meter("go-toolkit")
	v := expvar.Get("go-toolkit.http.client.conn_pools")
	if v == nil {
		return
	}

	var info pooledTransportPoolInfo
	if err := json.Unmarshal([]byte(v.String()), &info); err != nil {
		return
	}

	if _, err := meter.Int64ObservableGauge("go-toolkit.http.client.conn_pool",
		metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
			for pool, v := range info {
				for network, conns := range v {
					attr := []attribute.KeyValue{
						{
							Key:   "pool",
							Value: attribute.StringValue(pool),
						},
						{
							Key:   "network",
							Value: attribute.StringValue(network),
						},
					}
					observer.Observe(conns, metric.WithAttributes(attr...))
				}
			}
			return nil
		}),
	); err != nil {
		return
	}
}
