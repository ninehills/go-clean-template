package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// 自定义的 Logger，封装了zerolog Logger
type Logger struct {
	logger *zerolog.Logger
}

// 注入 request id
func (l *Logger) WithContext(c context.Context) *Logger {
	// 尝试从context中获取 *gin.Context
	value := c.Value(gin.ContextKey)
	if value == nil {
		// 如果并不包含 gin.Context，那么就不做任何处理
		return l
	} else {
		d := value.(*gin.Context)
		r := l.logger.With().Str("request_id", requestid.Get(d)).Logger()
		return &Logger{logger: &r}
	}
}

func (l *Logger) Log() *zerolog.Logger {
	return l.logger
}

// 如下方法是为了保持使用的简单
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

func New(level string, format string) (*Logger, error) {
	var out io.Writer
	if format == "text" {
		// 文本格式输出，参考 https://github.com/rs/zerolog#pretty-logging
		out = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	} else {
		// JSON 格式输出
		out = os.Stderr
	}

	// 设置日志级别
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return &Logger{}, fmt.Errorf("invalid log level: %s", level)
	}
	zerolog.SetGlobalLevel(logLevel)

	// 注册主机名
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Can't get hostname: %s", err))
	}

	// 使用短文件路径
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	// 初始化 Logger
	logger := zerolog.
		New(out).
		With().
		Caller().
		Timestamp().
		// 节点主机名，用于区分不同实例的日志
		Str("hostname", hostname).
		Logger()
	return &Logger{logger: &logger}, nil
}
