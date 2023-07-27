package replicate

import (
	"context"
	"net/http"
	"net/http/httptest"
)

type handler func(w http.ResponseWriter, r *http.Request)

type endpoint struct {
	path    string
	handler handler
}

func mockServer(endpoints ...endpoint) (string, func()) {
	mux := http.NewServeMux()
	for _, endpoint := range endpoints {
		mux.HandleFunc(endpoint.path, endpoint.handler)
	}

	server := httptest.NewServer(mux)
	return server.URL + "/%s", server.Close
}

func canceledContext() context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	return ctx
}
