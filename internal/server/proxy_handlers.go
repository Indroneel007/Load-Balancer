package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Next returns the next index in a round-robin fashion, thread-safe.
func (lb *LoadBalancer) Next(n int) int {
	if n <= 0 {
		return 0
	}
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()
	idx := lb.Current
	lb.Current = (lb.Current + 1) % n
	return idx
}

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

// ProxyRequestHandler returns an http handler that selects the next backend
// using the provided LoadBalancer and destinations slice on each request.
func ProxyRequestHandler(lb *LoadBalancer, destinations []string, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// choose next backend
		responseTimes := make(map[string]time.Duration)

		if len(responseTimes) == 0 {
			idx := lb.Next(len(destinations))
			targetStr := destinations[idx]
			targetURL, err := url.Parse(targetStr)
			if err != nil {
				http.Error(w, "bad backend url", http.StatusInternalServerError)
				return
			}

			fmt.Printf("Received request for %s, forwarding to %s\n", endpoint, targetURL.String())

			// Update the headers to allow for SSL redirection
			r.URL.Host = targetURL.Host
			r.URL.Scheme = targetURL.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = targetURL.Host

			// trim the endpoint prefix from the path
			path := r.URL.Path
			r.URL.Path = strings.TrimPrefix(path, endpoint)

			fmt.Printf("[ TinyRP ] Redirecting request to %s at %s\n", r.URL, time.Now().UTC())

			// Start a timer to measure response time
			startTime := time.Now()
			proxy := NewProxy(targetURL)
			proxy.ServeHTTP(w, r)
			responseTime := time.Since(startTime)

			responseTimes[targetStr] = responseTime
		} else {
			// Find the backend with the lowest response time
			minTime := time.Duration(0)
			minBackend := ""
			for backend, time := range responseTimes {
				if time < minTime || minTime == 0 {
					minTime = time
					minBackend = backend
				}
			}

			// Choose the backend with the lowest response time
			targetStr := minBackend
			targetURL, err := url.Parse(targetStr)
			if err != nil {
				http.Error(w, "bad backend url", http.StatusInternalServerError)
				return
			}

			// Update the headers to allow for SSL redirection
			r.URL.Host = targetURL.Host
			r.URL.Scheme = targetURL.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = targetURL.Host

			// trim the endpoint prefix from the path
			path := r.URL.Path
			r.URL.Path = strings.TrimPrefix(path, endpoint)

			// Start a timer to measure response time
			startTime := time.Now()
			proxy := NewProxy(targetURL)
			proxy.ServeHTTP(w, r)
			responseTime := time.Since(startTime)

			// Store the response time for this backend
			responseTimes[targetStr] = responseTime
		}
	}
}
