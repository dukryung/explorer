package log

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hessegg/nikto-explorer/types/config"
)

type Logger struct {
	prefix string
	enable bool
	level  int
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	error  *log.Logger
}

func (l *Logger) write(logger *log.Logger, logLv int, v ...interface{}) {
	if !l.enable {
		return
	}

	if l.level > logLv {
		return
	}

	prefix := fmt.Sprintf(color.CyanString("x") + "=" + l.prefix + " ")

	logger.Output(3, prefix+fmt.Sprint(v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.write(l.debug, 2, v)
}

func (l *Logger) Info(v ...interface{}) {
	l.write(l.info, 3, v)
}

func (l *Logger) Warn(v ...interface{}) {
	l.write(l.warn, 4, v)
}

func (l *Logger) Err(err error) error {
	l.write(l.error, 5, fmt.Sprint(err))
	return err
}

func (l *Logger) Error(v ...interface{}) {
	l.write(l.error, 5, v)
}

func NewLogger(prefix string, config config.LogConfig) *Logger {
	prefix = strings.Trim(prefix, " ")

	flag := log.Ltime | log.Lshortfile
	l := Logger{
		prefix: prefix,
		enable: config.Enable,
		level:  config.Level,
		debug:  log.New(os.Stdout, color.BlueString("DEBUG "), flag),
		info:   log.New(os.Stdout, color.CyanString("INFO "), flag),
		warn:   log.New(os.Stdout, color.YellowString("WARN "), flag),
		error:  log.New(os.Stdout, color.RedString("ERROR "), flag),
	}

	return &l
}
