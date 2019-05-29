/*
Copyright (c) JSC iCore.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
*/

package routegroup_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/i-core/routegroup"
)

func TestRouter(t *testing.T) {
	var routes testRouteAdder = []*testRoute{
		{method: http.MethodGet, path: "/foo", handler: &testHandler{}},
		{method: http.MethodGet, path: "/bar", handler: &testHandler{}},
	}

	router := routegroup.NewRouter()
	router.AddRoutes(routes, "")

	for _, route := range routes {
		r := httptest.NewRequest(route.method, route.path, nil)
		router.ServeHTTP(httptest.NewRecorder(), r)
		if !route.handler.called {
			t.Errorf("handler for route %s %s is not called", route.method, route.path)
		}
	}
}

func TestRoutesWithPrefix(t *testing.T) {
	var routes testRouteAdder = []*testRoute{
		{method: http.MethodGet, path: "/foo", handler: &testHandler{}},
		{method: http.MethodGet, path: "/bar", handler: &testHandler{}},
	}

	router := routegroup.NewRouter()
	router.AddRoutes(routes, "/prefix")

	for _, route := range routes {
		r := httptest.NewRequest(route.method, "/prefix"+route.path, nil)
		router.ServeHTTP(httptest.NewRecorder(), r)
		if !route.handler.called {
			t.Errorf("handler for route %s %s is not called", route.method, route.path)
		}
	}
}

func TestRouterWithMiddleware(t *testing.T) {
	var routes testRouteAdder = []*testRoute{
		{method: http.MethodGet, path: "/foo", handler: &testHandler{}},
		{method: http.MethodGet, path: "/bar", handler: &testHandler{}},
	}

	mw := &testMiddleware{}
	router := routegroup.NewRouter(mw.newInstance())
	router.AddRoutes(routes, "/prefix")

	for _, route := range routes {
		r := httptest.NewRequest(route.method, "/prefix"+route.path, nil)
		router.ServeHTTP(httptest.NewRecorder(), r)
		if !route.handler.called {
			t.Errorf("handler for route %s %s is not called", route.method, route.path)
		}
	}

	if wantCalls := len(routes); mw.calls != wantCalls {
		t.Errorf("got %d call(s) of middleware; want %d call(s)", mw.calls, wantCalls)
	}
}

func TestPathParam(t *testing.T) {
	var (
		method    = http.MethodGet
		path      = "/foo/:param"
		url       = "/foo/bar"
		wantParam = "bar"
		gotParam  string
	)

	handler := &testHandler{
		fn: func(w http.ResponseWriter, r *http.Request) {
			gotParam = routegroup.PathParam(r.Context(), "param")
		},
	}
	router := routegroup.NewRouter()
	router.AddRoutes(testRouteAdder([]*testRoute{{method: method, path: path, handler: handler}}), "")
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(method, url, nil))

	if gotParam != wantParam {
		t.Errorf("got params %q; want param %q", gotParam, wantParam)
	}
}

type testRoute struct {
	method  string
	path    string
	handler *testHandler
}

type testRouteAdder []*testRoute

func (a testRouteAdder) AddRoutes(apply func(m, p string, h http.Handler, mws ...func(http.Handler) http.Handler)) {
	for _, r := range a {
		apply(r.method, r.path, r.handler)
	}
}

type testHandler struct {
	called bool
	fn     func(http.ResponseWriter, *http.Request)
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	if h.fn != nil {
		h.fn(w, r)
	}
}

type testMiddleware struct {
	calls int
}

func (m *testMiddleware) newInstance() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.calls++
			h.ServeHTTP(w, r)
		})
	}
}
