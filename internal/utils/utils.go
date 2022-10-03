package utils

import (
	"bufio"
	"os"
	"strings"
	"synchfolder/internal/logger"
)

type Config struct {
}

var ConfigPath string

var logError logger.LogMessage = logger.LogMessage{LogType: logger.LogCritical, Ref: "", Message: ""}

func GetConfig() (map[string]string, error) {

	logError.Ref = "GetConfig"

	result := map[string]string{}

	file, err := os.Open(ConfigPath)

	if err != nil {

		logError.Message = "error reading config.txt : " + err.Error()
		logger.LogChan <- logError

		return nil, err

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		str := scanner.Text()
		if strings.Contains(str, "=") {
			tmp := strings.Split(str, "=")
			tmp[1] = strings.TrimSpace(tmp[1])
			result[tmp[0]] = tmp[1]
		}

	}

	return result, nil
}
