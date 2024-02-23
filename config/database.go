package config

import "fmt"

// Database represents the configuration details for a database connection.
type Database struct {
	Host          string
	Port          string
	User          string
	Password      string
	Database      string
	MaxConnection int32
}

// Address returns the formatted string for the database connection address.
func (e *Database) Address() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		e.Host,
		e.Port,
		e.User,
		e.Password,
		e.Database,
	)
}
