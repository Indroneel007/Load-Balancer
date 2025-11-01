package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

// ProxyRequestHandler returns an http handler that selects the next backend
// using the provided LoadBalancer and destinations slice on each request.
func ProxyRequestHandler(lb *LoadBalancer, destinations []string, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// choose next backend
		idx := 0
		if len(destinations) > 1 {
			idx = lb.Next(len(destinations))
		}
		targetStr := destinations[idx]
		targetURL, err := url.Parse(targetStr)
		if err != nil {
			http.Error(w, "bad backend url", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Received request for %s, forwarding to %s\n", endpoint, targetURL.String())

		proxy := NewProxy(targetURL)

		// Update the headers to allow for SSL redirection
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = targetURL.Host

		// trim the endpoint prefix from the path
		path := r.URL.Path
		r.URL.Path = strings.TrimPrefix(path, endpoint)

		// Note that ServeHTTP is non blocking and uses a goroutine under the hood
		fmt.Printf("[ TinyRP ] Redirecting request to %s at %s\n", r.URL, time.Now().UTC())
		proxy.ServeHTTP(w, r)
	}
}
