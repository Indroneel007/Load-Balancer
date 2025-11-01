package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Indroneel007/Load-Balancer/internal/config"
)

type LoadBalancer struct {
	Current int
	Mutex   sync.Mutex
}

func Run() error {
	configs, err := config.NewConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create a new Router
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", ping)

	// Iterating through the configuration resource and registering them
	// into the router. For each resource we create a LoadBalancer instance
	// and give the handler access to the list of backend URLs. The handler
	// will pick the next backend on each request (round-robin).
	for _, resource := range configs.Resources {
		// Normalize destinations: prefer Destinations slice, but fall back to single Destination_Url
		dests := resource.Destinations
		if len(dests) == 0 && resource.Destination_Url != "" {
			dests = []string{resource.Destination_Url}
		}
		if len(dests) == 0 {
			fmt.Printf("resource %s has no destinations, skipping\n", resource.Name)
			continue
		}
		fmt.Printf("resource %s destinations: %v\n", resource.Name, dests)
		lb := &LoadBalancer{}
		mux.HandleFunc(resource.Endpoint, ProxyRequestHandler(lb, dests, resource.Endpoint))
	}
	// Running proxy server
	addr := configs.Server.Host + ":" + configs.Server.Listen_port
	fmt.Printf("listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}
	return nil
}
