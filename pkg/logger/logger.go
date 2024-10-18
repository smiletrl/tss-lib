package logger

type Logger interface {
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}
