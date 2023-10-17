package httprouter

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

// Handler is a type that handles an http request within our framework.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Middleware is a function designed to run some code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given Handler.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// wrapMiddleware creates a new handler by wrapping mw around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(handler http.HandlerFunc, mw []Middleware) http.HandlerFunc {
	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

// Config allows configuring a Router instance.
type Config struct {
	ErrorHandlerFunc            ErrorHandlerFunc
	NotFoundHandler             http.Handler
	HealthCheckLivenessHandler  http.Handler
	HealthCheckReadinessHandler http.Handler
	EnableProfiler              bool

	Mw []Middleware
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes.
type Router struct {
	mux            *chi.Mux
	mw             []Middleware
	errHandlerFunc ErrorHandlerFunc
}

// New instantiates a `Router` with the given configuration.
func New(cfg Config) *Router {
	mux := chi.NewRouter()

	if cfg.NotFoundHandler != nil {
		mux.NotFound(cfg.NotFoundHandler.ServeHTTP)
	}

	if cfg.HealthCheckLivenessHandler != nil {
		mux.Get("/liveness", cfg.HealthCheckLivenessHandler.ServeHTTP)
	}

	if cfg.HealthCheckReadinessHandler != nil {
		mux.Get("/readiness", cfg.HealthCheckReadinessHandler.ServeHTTP)
	}

	if cfg.EnableProfiler {
		mux.Mount("/debug", middleware.Profiler())
	}

	return &Router{
		mux:            mux,
		mw:             cfg.Mw,
		errHandlerFunc: cfg.ErrorHandlerFunc,
	}
}

// Group creates a new RouteGroup with the given p prefix and middlewares which are
// chained after this Router's middlewares.
func (r *Router) Group(p string, mw ...Middleware) *RouteGroup {
	return &RouteGroup{
		router: r,
		path:   p,
		mw:     mw,
	}
}

// Method adds the route `pattern` that matches `method` http method to
// execute the `handler` http.Handler wrapped by `mw`.
func (r *Router) Method(method, pattern string, handler Handler, mw ...Middleware) {
	r.mux.Method(method, pattern, r.wrapHandler(r.handlerAdapter(handler), mw...))
}

// Any adds the route `pattern` that matches any http method to execute the `handler` http.Handler wrapped by `mw`.
func (r *Router) Any(pattern string, handler Handler, mw ...Middleware) {
	r.mux.Handle(pattern, r.wrapHandler(r.handlerAdapter(handler), mw...))
}

func (r *Router) handlerAdapter(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := handler(w, req)
		if err == nil {
			return
		}

		var hErr HandlerError
		if r.errHandlerFunc != nil {
			hErr = r.errHandlerFunc(err, DefaultHandlerError)
		} else {
			hErr = DefaultHandlerError(err)
		}

		if hErr.Notify {
			span := trace.SpanFromContext(req.Context())
			defer span.End()

			notifyErr(req, err, hErr.StatusCode)
		}

		_ = RespondJSON(w, hErr.StatusCode, hErr.Error)
	}
}

func (r *Router) wrapHandler(handler http.HandlerFunc, mw ...Middleware) http.HandlerFunc {
	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(handler, mw)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(handler, r.mw)

	return func(w http.ResponseWriter, r *http.Request) {
		// Populate context with URI params for later retrieval.
		routeCtx := chi.RouteContext(r.Context())
		routeParams := routeCtx.URLParams

		params := make(URIParams, len(routeParams.Keys))
		for i := range routeParams.Keys {
			params[routeParams.Keys[i]] = routeParams.Values[i]
		}

		r = r.WithContext(WithParams(r.Context(), params))
		handler(w, r)
	}
}

// Get is a shortcut for r.Method(http.MethodGet, pattern, handle, mw).
func (r *Router) Get(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodGet, pattern, handler, mw...)
}

// Head is a shortcut for r.Method(http.MethodHead, pattern, handle, mw).
func (r *Router) Head(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodHead, pattern, handler, mw...)
}

// Options is a shortcut for r.Method(http.MethodOptions, pattern, handle, mw).
func (r *Router) Options(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodOptions, pattern, handler, mw...)
}

// Post is a shortcut for r.Method(http.MethodPost, pattern, handle, mw).
func (r *Router) Post(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodPost, pattern, handler, mw...)
}

// Put is a shortcut for r.Method(http.MethodPut, pattern, handle, mw).
func (r *Router) Put(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodPut, pattern, handler, mw...)
}

// Patch is a shortcut for r.Method(http.MethodPatch, pattern, handle, mw).
func (r *Router) Patch(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodPatch, pattern, handler, mw...)
}

// Delete is a shortcut for r.Method(http.MethodDelete, pattern, handle, mw).
func (r *Router) Delete(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodDelete, pattern, handler, mw...)
}

// Trace is a shortcut for r.Method(http.MethodTrace, pattern, handle, mw).
func (r *Router) Trace(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodTrace, pattern, handler, mw...)
}

// ServeHTTP conforms to the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Route describes the details of a routing handler.
type Route struct {
	Method      string
	Route       string
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

// Routes returns the routing tree in an easily traversable structure.
func (r *Router) Routes() ([]Route, error) {
	var routes []Route
	walkFunc := func(method string, route string, handler http.Handler, mw ...func(http.Handler) http.Handler) error {
		routes = append(routes, Route{
			Method:      method,
			Route:       route,
			Handler:     handler,
			Middlewares: mw,
		})
		return nil
	}

	if err := chi.Walk(r.mux, walkFunc); err != nil {
		return nil, fmt.Errorf("generating routes %w", err)
	}

	return routes, nil
}
