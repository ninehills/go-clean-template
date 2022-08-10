package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
)

type Logger interface {
	logur.Logger

	// Format strings
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	// WithFields returns a new logger with fields added to the current logger.
	WithFields(fields map[string]interface{}) Logger

	// WithField is a shortcut for WithFields(map[string]interface{}{key: value}).
	WithField(key string, value interface{}) Logger

	// WithContext = Ctx 从context中读取request_id
	Ctx(ctx context.Context) Logger

	// WithError = Err returns a new logger with the given error added to the current logger.
	Err(err error) Logger
}

// Config holds details necessary for logging.
type Config struct {
	// Format specifies the output log format.
	// Accepted values are: json, text(default).
	Format string

	// Level is the minimum log level that should appear on the output.
	Level string

	// NoColor makes sure that no log output gets colorized.
	NoColor bool
}

type logger struct {
	logger logur.Logger
}

// Trace implements the logur.Logger interface.
func (l *logger) Trace(msg string, fields ...map[string]interface{}) {
	l.logger.Trace(msg, fields...)
}

// Tracef -.
func (l *logger) Tracef(format string, args ...interface{}) {
	l.Trace(fmt.Sprintf(format, args...))
}

// Debug implements the logur.Logger interface.
func (l *logger) Debug(msg string, fields ...map[string]interface{}) {
	l.logger.Debug(msg, fields...)
}

// Debugf -.
func (l *logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Info implements the logur.Logger interface.
func (l *logger) Info(msg string, fields ...map[string]interface{}) {
	l.logger.Info(msg, fields...)
}

// Infof -.
func (l *logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn implements the logur.Logger interface.
func (l *logger) Warn(msg string, fields ...map[string]interface{}) {
	l.logger.Warn(msg, fields...)
}

// Warnf -.
func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Error implements the logur.Logger interface.
func (l *logger) Error(msg string, fields ...map[string]interface{}) {
	l.logger.Error(msg, fields...)
}

// Errorf -.
func (l *logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// WithFields returns a new logger instance that attaches the given fields to every subsequent log call.
func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &logger{
		logger: logur.WithFields(l.logger, fields),
	}
}

// WithField is a shortcut for WithFields(logger, map[string]interface{}{key: value}).
func (l *logger) WithField(key string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{key: value})
}

// Ctx 从context中读取request_id.
func (l *logger) Ctx(c context.Context) Logger {
	// 尝试从context中获取 *gin.Context
	value := c.Value(gin.ContextKey)
	if value == nil {
		// 如果并不包含 gin.Context，那么就不做任何处理
		return l
	}

	d, ok := value.(*gin.Context)
	if !ok {
		// 如果不是 *gin.Context，那么就不做任何处理
		return l
	}

	return l.WithField("request_id", requestid.Get(d))
}

// WithError = Err.
func (l *logger) Err(err error) Logger {
	return l.WithField("error", err)
}

// New creates a new logger.
func New(config Config) Logger {
	l := logrus.New()

	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.TextFormatter{
		DisableColors:             config.NoColor,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
	})
	// 因为嵌套了一层，caller 无法获取真正的调用者，此处注释，后续寻找方法
	// 参见 https://github.com/sirupsen/logrus/pull/989
	// l.SetReportCaller(true)
	switch config.Format {
	case "text":
		// Already the default

	case "json":
		l.SetFormatter(&logrus.JSONFormatter{})
	}

	if level, err := logrus.ParseLevel(config.Level); err == nil {
		l.SetLevel(level)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Can't get hostname: %s", err))
	}

	ll := &logger{
		logger: logrusadapter.New(l),
	}

	return ll.WithField("hostname", hostname)
}

// SetStandardLogger sets the global logger's output to a custom logger instance.
func SetStandardLogger(l Logger) {
	log.SetOutput(logur.NewLevelWriter(l, logur.Info))
}
