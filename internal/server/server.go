package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Indroneel007/Load-Balancer/internal/config"
)

func Run() error {
	configs, err := config.NewConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create a new Router
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", ping)

	// Iterating through the configuration resource and registering them
	// into the router.
	for _, resource := range configs.Resources {
		url, _ := url.Parse(resource.Destination_Url)
		fmt.Printf("parsed url is %v\n", url)
		proxy := NewProxy(url)
		mux.HandleFunc(resource.Endpoint, ProxyRequestHandler(proxy, url, resource.Endpoint))
	}
	// Running proxy server
	addr := configs.Server.Host + ":" + configs.Server.Listen_port
	fmt.Printf("listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}
	return nil
}
