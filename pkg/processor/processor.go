package processor

import "context"

// Processor represents the interface for a process runtime
type Processor interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Factory represents the connect factor like db, cache,eg...
type Factory interface {
	Connect(context.Context) error
	Close(context.Context) error
}
