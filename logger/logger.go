package logger

import (
	"fmt"
	"os"
	"time"
)

const Error string = "ERR "
const Warn string = "WARN"
const Info string = "INFO"

type Logger interface {
	Log(kind string, prefix string, fmt string, args ...any) error
}

type FileLogger struct {
	file *os.File
}

type TerminalLogger struct {
	FileLogger
	SaveToFile bool
}

func (logger FileLogger) Log(kind string, prefix string, format string, args ...any) error {
	formatted := fmt.Sprintf(format, args...)
	data := fmt.Sprintf("[%s %s] [%s]: %s\r\n", time.Now().Local().Format(time.UnixDate), kind, prefix, formatted)

	_, err := logger.file.WriteString(data)

	return err
}

func (logger TerminalLogger) Log(kind string, prefix string, format string, args ...any) error {
	var err error

	formatted := fmt.Sprintf(format, args...)
	data := fmt.Sprintf("[%s %s] [%s]: %s\r\n", time.Now().Local().Format(time.UnixDate), kind, prefix, formatted)

	_, err = fmt.Print(data)

	if err != nil {
		return err
	}

	if logger.SaveToFile && logger.file != nil {
		_, err = logger.file.WriteString(data)
	}

	return err
}
