package log

import (
	"fmt"
	"log"
	"os"
)

var (
	_log *myLogger
)

func Init(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 777)
	if err != nil {
		return err
	}
	logger := log.New(file, "[dcss]", log.LstdFlags|log.Lshortfile)
	_log = &myLogger{logger: logger}
	return nil
}

type myLogger struct {
	logger *log.Logger
}

func Debug(data interface{}) {
	_log.logger.Printf("[DEBUG] %s", data)
}

func Info(data interface{}) {
	_log.logger.Printf("[INFO] %s", data)
}

func Error(data interface{}) {
	_log.logger.Printf("[ERROR] %s", data)
}

func DebugF(format string, data ...interface{}) {
	output := fmt.Sprintf(format, data...)
	_log.logger.Printf("[DEBUG] %s", output)
}

func InfoF(format string, data ...interface{}) {
	output := fmt.Sprintf(format, data...)
	_log.logger.Printf("[INFO] %s", output)
}

func ErrorF(format string, data ...interface{}) {
	output := fmt.Sprintf(format, data...)
	_log.logger.Printf("[ERROR] %s", output)
}
