package log

// Logger - simplified logger abstraction
type Logger interface {
	Debug(...any)
	Info(...any)
	Error(error, ...any)
}

// MainLogger - Logger which can also report fatal errors
type MainLogger interface {
	Logger
	FatalError(error, ...any)
}
