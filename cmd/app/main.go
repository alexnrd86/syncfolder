package main

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"synchfolder/internal/logger"
	"synchfolder/internal/synch"
	"synchfolder/internal/utils"
	"time"
)

func main() {

	var root string //root path

	var cfgMap map[string]string //map for config storing

	var wgMain sync.WaitGroup

	logInfo := logger.LogMessage{LogType: logger.LogInfo, Ref: "main", Message: "start"} //struct for log message storing type INFO

	//initialising paths

	root, _ = filepath.Abs("./")
	logger.LogPath = root + "/logs/log.txt"
	utils.ConfigPath = root + "/configs/config.txt"

	ctxLogger, cancelLogger := context.WithCancel(context.Background())

	// starting logger
	go func() {
		logger.Logger(ctxLogger)
	}()

	defer func() {
		cancelLogger()
	}()

	cfgMap, err := utils.GetConfig() //read config.txt

	if err != nil {
		fmt.Println("Error reading config.txt. App terminated")
		return
	}

	if err = logger.SetLogLevel(cfgMap["loglevel"]); err != nil {
		fmt.Println(err.Error())
	}

	logger.LogChan <- logInfo //log app start

L:
	for {

		select {

		case <-synch.CriticalChan:

			fmt.Println("application is closed unexpectedly")

			break L

		default:
			wgMain.Add(1)

			go func() {
				defer wgMain.Done()
				_ = synch.CheckMasterFolder(cfgMap["sourcepath"], cfgMap["synchpath"])
			}()

			wgMain.Add(1)

			go func() {
				defer wgMain.Done()
				_ = synch.CheckSlaveFolder(cfgMap["sourcepath"], cfgMap["synchpath"])
			}()

			wgMain.Wait()

			time.Sleep(3 * time.Second)

		}
	}

}
