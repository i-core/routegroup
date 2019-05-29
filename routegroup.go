/*
Copyright (c) JSC iCore.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
*/

// Package routegroup provides a modular approach to organizing HTTP handlers
// and encapsulates an HTTP router's implementation.
package routegroup // import "github.com/i-core/routegroup"

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type pathParamsContextKeyType struct{}

var pathParamsContextKey = pathParamsContextKeyType{}

// RouteAdder is an interface that is used for registering routes.
type RouteAdder interface {
	AddRoutes(apply func(m, p string, h http.Handler, mws ...func(http.Handler) http.Handler))
}

// Router is a HTTP router that serves HTTP handlers.
type Router struct {
	http.Handler
	impl *httprouter.Router
}

// NewRouter creates a new Router.
func NewRouter(mws ...func(http.Handler) http.Handler) *Router {
	var ctrs []alice.Constructor
	for _, mw := range mws {
		ctrs = append(ctrs, mw)
	}
	impl := httprouter.New()
	return &Router{Handler: alice.New(ctrs...).Then(impl), impl: impl}
}

// AddRoutes adds RouteAdder's routes with an optional prefix.
func (rt *Router) AddRoutes(adder RouteAdder, prefix string) {
	adder.AddRoutes(func(m, p string, h http.Handler, mws ...func(http.Handler) http.Handler) {
		var ctrs []alice.Constructor
		for _, mw := range mws {
			ctrs = append(ctrs, mw)
		}
		rt.handler(m, prefix+p, alice.New(ctrs...).Then(h))
	})
}

func (rt *Router) handler(method, path string, handler http.Handler) {
	rt.impl.Handle(method, path,
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			v := make(map[string]string)
			for _, p := range ps {
				v[p.Key] = p.Value
			}
			ctx := context.WithValue(r.Context(), pathParamsContextKey, v)
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
		},
	)
}

// PathParam returns an URL path's parameter by a name.
func PathParam(ctx context.Context, key string) string {
	ps, ok := ctx.Value(pathParamsContextKey).(map[string]string)
	if !ok {
		return ""
	}
	return ps[key]
}
