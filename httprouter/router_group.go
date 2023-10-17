package httprouter

import (
	"net/http"
	"path"
)

// RouteGroup represents a group of routes that share the same path prefix and middlewares.
type RouteGroup struct {
	router *Router
	path   string
	mw     []Middleware
}

// Group creates a new RouteGroup with the given path relative to the existing RouteGroup path
// and middlewares which are chained after this RouteGroup's middlewares.
func (g *RouteGroup) Group(p string, mw ...Middleware) *RouteGroup {
	g.mw = append(g.mw, mw...)
	return g.router.Group(path.Join(g.path, p), g.mw...)
}

// Method adds the route `pattern` that matches `method` http method to
// execute the `handler` http.Handler wrapped by `mw`.
func (g *RouteGroup) Method(method, pattern string, handler Handler, mw ...Middleware) {
	mw = append(mw, g.mw...)
	g.router.Method(method, path.Join(g.path, pattern), handler, mw...)
}

// Any adds the route `pattern` that matches any http method to execute the `handler` http.Handler wrapped by `mw`.
func (g *RouteGroup) Any(pattern string, handler Handler, mw ...Middleware) {
	mw = append(mw, g.mw...)
	g.router.Any(path.Join(g.path, pattern), handler, mw...)
}

// Get is a shortcut for g.Method(http.MethodGet, pattern, handle, mw).
func (g *RouteGroup) Get(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodGet, pattern, handler, mw...)
}

// Head is a shortcut for g.Method(http.MethodHead, pattern, handle, mw).
func (g *RouteGroup) Head(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodHead, pattern, handler, mw...)
}

// Options is a shortcut for g.Method(http.MethodOptions, pattern, handle, mw).
func (g *RouteGroup) Options(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodOptions, pattern, handler, mw...)
}

// Post is a shortcut for g.Method(http.MethodPost, pattern, handle, mw).
func (g *RouteGroup) Post(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodPost, pattern, handler, mw...)
}

// Put is a shortcut for g.Method(http.MethodPut, pattern, handle, mw).
func (g *RouteGroup) Put(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodPut, pattern, handler, mw...)
}

// Patch is a shortcut for g.Method(http.MethodPatch, pattern, handle, mw).
func (g *RouteGroup) Patch(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodPatch, pattern, handler, mw...)
}

// Delete is a shortcut for g.Method(http.MethodDelete, pattern, handle, mw).
func (g *RouteGroup) Delete(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodDelete, pattern, handler, mw...)
}

// Trace is a shortcut for g.Method(http.MethodTrace, pattern, handle, mw).
func (g *RouteGroup) Trace(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodTrace, pattern, handler, mw...)
}
