package logger

import (
	"bufio"
	"context"
	"errors"
	"os"
	"strings"
	"time"
)

const (
	LogError    string = "ERROR"
	LogInfo     string = "INFO"
	LogCritical string = "CRITICAL"
)

type LogMessage struct {
	LogType string
	Ref     string
	Message string
}

var LogChan chan LogMessage

var LogPath, LogLevel string

func init() {
	LogChan = make(chan LogMessage, 100)
	LogLevel = LogError

}

func log(logtype, ref, message string) error {

	logfile, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer logfile.Close()

	t := time.Now()

	str := t.Format("02-01-2006 15:04:05") + " - " + logtype + " - FUNC: " + ref + "; LOG: " + message + ";\n"

	w := bufio.NewWriter(logfile)

	_, err = w.WriteString(str)

	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

func Logger(ctx context.Context) {
	logfile, _ := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logfile.Close()

	defer close(LogChan)

	for {
		select {
		case message := <-LogChan:

			switch {
			case LogLevel == LogInfo:
				_ = log(message.LogType, message.Ref, message.Message)

			case LogLevel == LogError:
				if message.LogType != LogInfo {
					_ = log(message.LogType, message.Ref, message.Message)

				}
			case LogLevel == LogCritical:
				if message.LogType == LogCritical {
					_ = log(message.LogType, message.Ref, message.Message)
				}

			}

		case <-ctx.Done():

			return
		}
	}
}

func SetLogLevel(logLevel string) error {

	logLevel = strings.ToUpper(logLevel)

	if logLevel == LogInfo || logLevel == LogError || logLevel == LogCritical {
		LogLevel = logLevel

	} else {
		return errors.New("log level is not set in config. Default log level ERROR")
	}

	return nil

}
