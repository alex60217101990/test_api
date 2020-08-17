package logger

type Logger interface {
	GetNativeLogger() interface{}
	Info(msg string)
	Infof(tpl string, args map[string]interface{})
	Error(err error)
	Errorf(tpl string, args map[string]interface{})
	Warn(msg string)
	Warnf(tpl string, args map[string]interface{})
	Fatal(err error)
	Fatalf(tpl string, args map[string]interface{})
	Printf(format string, args ...interface{})
	Close()
}
