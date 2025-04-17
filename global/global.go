package global

// Logger is an interface for log messages
type Logger interface {
	Debug(string)
	Info(string)
	Notice(string)
	Warning(string)
	Error(string)
	Fatal(string)
	Debugf(string, ...any)
	Infof(string, ...any)
	Noticef(string, ...any)
	Warningf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
	Close()
}
