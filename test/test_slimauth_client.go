package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/cmstar/go-logx"
	"github.com/cmstar/go-webapi"
	client "github.com/cmstar/go-webapi-client"
	"github.com/cmstar/go-webapi-client/slimauth_client"
	"github.com/cmstar/go-webapi/slimauth"
)

const (
	_KEY    = "my-app"
	_SECRET = "my-secret"
)

func main() {
	ts := newTestServer()
	runClient(ts.URL)
}

func newTestServer() *httptest.Server {
	handler := slimauth.NewSlimAuthApiHandler(slimauth.SlimAuthApiHandlerOption{
		SecretFinder: func(accessKey string) string {
			if accessKey == _KEY {
				return _SECRET
			}
			return ""
		},
	})

	handler.RegisterMethod(webapi.ApiMethod{
		Name: "Test",
		Value: reflect.ValueOf(func(req struct{ S1, S2 string }) string {
			return req.S1 + "\n" + req.S2
		}),
	})

	logger := logx.NopLogger
	handlerFunc := webapi.CreateHandlerFunc(handler, logx.NewSingleLoggerLogFinder(logger))
	ts := httptest.NewServer(http.HandlerFunc(handlerFunc))
	return ts
}

func runClient(uri string) {
	c := slimauth_client.NewClient()
	c.SetConfig(map[string]any{
		"Key":    _KEY,
		"Secret": _SECRET,
		"Uri":    uri + "?Test",
		"Param":  `{"S1": 123, "S2": "line1\nline2"}`,
	})

	op := &client.RunOption{
		Clients: []client.Client{c},
	}
	client.Run(op)
}
