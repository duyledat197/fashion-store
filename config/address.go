package config

import "fmt"

type Endpoint struct {
	Host string
	Port string
}

func (e *Endpoint) Address() string {
	return fmt.Sprintf("%s:%s", e.Host, e.Port)
}
