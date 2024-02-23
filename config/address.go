package config

import "fmt"

// Endpoint represents the configuration for a network endpoint.
type Endpoint struct {
	Host string // Host is the IP address or domain name of the endpoint.
	Port string // Port is the network port number of the endpoint.
}

// Address returns the formatted address string combining the Host and Port.
func (e *Endpoint) Address() string {
	return fmt.Sprintf("%s:%s", e.Host, e.Port)
}
