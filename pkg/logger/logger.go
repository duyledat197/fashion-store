package logger

// Logger is a presentation for logger engine apply logger pattern which help logging for system.
type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Info(msg string, args ...any)
}
