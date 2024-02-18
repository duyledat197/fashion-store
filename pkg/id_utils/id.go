package id_utils

// IDGenerator is an exporter for common interface for various id generator.
type IDGenerator interface {
	String() string
	Int64() int64
}
