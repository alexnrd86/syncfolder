package synch

import (
	"errors"
	"io"
	"os"
	"sync"
	"synchfolder/internal/logger"
)

var CriticalChan chan struct{}

func init() {

	CriticalChan = make(chan struct{}, 2)
}

// func that check master folder and run a goroutine for every file in the folder.
// If it finds a subfolder it runs itself for subfolder
func CheckMasterFolder(masterPath, slavePath string) error {

	var logCritical logger.LogMessage = logger.LogMessage{LogType: logger.LogCritical, Ref: "CheckMasterFolder", Message: ""}

	folder, err := os.ReadDir(masterPath)

	if err != nil {
		logCritical.Message = "error reading master folder: " + err.Error()
		logger.LogChan <- logCritical
		CriticalChan <- struct{}{}
		return err
	}

	var wgCMF sync.WaitGroup

	for _, entry := range folder {

		if entry.IsDir() {

			dirInfo, _ := entry.Info()

			err = checkFolder(entry.Name(), slavePath, dirInfo.Mode().Perm())

			if err == nil {
				_ = CheckMasterFolder(masterPath+"/"+entry.Name(), slavePath+"/"+entry.Name())
			}

		} else {

			wgCMF.Add(1)
			go func() {

				defer wgCMF.Done()
				_ = checkFile(entry, masterPath, slavePath)

			}()
		}

	}

	wgCMF.Wait()

	return nil

}

// func check if the file exists in the slave folder. If not - copy file from source folder
func checkFile(entry os.DirEntry, masterPath string, slavePath string) error {

	var logError logger.LogMessage = logger.LogMessage{LogType: logger.LogError, Ref: "checkFile", Message: ""}

	if entry.IsDir() {
		return errors.New("not a file")
	}

	folder, err := os.ReadDir(slavePath)

	if err != nil {
		logError.Message = err.Error()
		logger.LogChan <- logError
		return err

	}

	exist := false

	for _, slEntry := range folder {

		if slEntry.Name() == entry.Name() {

			exist = true

			msFileInfo, _ := entry.Info()
			slFileInfo, _ := slEntry.Info()

			if msFileInfo.Size() == slFileInfo.Size() {
				return nil

			} else {
				err := os.Remove(slavePath + "/" + slFileInfo.Name())

				if err != nil {

					logError.Message = err.Error()
					logger.LogChan <- logError
					return err

				}

				err = copyFile(masterPath+"/"+msFileInfo.Name(), slavePath+"/"+slFileInfo.Name())

				if err != nil {

					logError.Message = err.Error()
					logger.LogChan <- logError

				}

			}
		}
	}

	if !exist {

		err = copyFile(masterPath+"/"+entry.Name(), slavePath+"/"+entry.Name())

		if err != nil {

			logError.Message = err.Error()
			logger.LogChan <- logError

		}
	}

	return nil

}

// func that copy file from inPath to outPath
func copyFile(inPath, outPath string) error {

	var logInfo logger.LogMessage = logger.LogMessage{LogType: logger.LogInfo, Ref: "copyFile", Message: ""}

	in, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	logInfo.Message = "Copy file " + inPath + " to " + outPath
	logger.LogChan <- logInfo

	return nil
}

// func check if a folder exists in slave folder. If not, create the folder in slave
func checkFolder(name, slavePath string, perm os.FileMode) error {

	var logInfo logger.LogMessage = logger.LogMessage{LogType: logger.LogInfo, Ref: "checkFolder", Message: ""}
	var logError logger.LogMessage = logger.LogMessage{LogType: logger.LogError, Ref: "checkFolder", Message: ""}

	folder, err := os.ReadDir(slavePath)

	if err != nil {

		logError.Message = err.Error()
		logger.LogChan <- logError

		return err

	}

	for _, slEntry := range folder {

		if slEntry.Name() == name {

			return nil
		}

	}

	err = os.Mkdir(slavePath+"/"+name, perm)

	if err != nil {

		logError.Message = err.Error()
		logger.LogChan <- logError

		return err

	}

	logInfo.Message = "Folder " + name + " created in " + slavePath
	logger.LogChan <- logInfo

	return nil
}

// func check slave folder and runs goroutine for every file to check if it still exists in the source folder
// if it finds a subfolder it runs itself for subfolder
func CheckSlaveFolder(masterPath, slavePath string) error {

	var logCritical logger.LogMessage = logger.LogMessage{LogType: logger.LogCritical, Ref: "CheckSlaveFolder", Message: ""}

	var wgCSF sync.WaitGroup

	folder, err := os.ReadDir(slavePath)

	if err != nil {
		logCritical.Message = err.Error()
		logger.LogChan <- logCritical
		CriticalChan <- struct{}{}
		return err
	}

	for _, entry := range folder {

		if entry.IsDir() {

			deleted, err := removeFolder(entry.Name(), masterPath, slavePath)

			if !deleted && err == nil {
				_ = CheckSlaveFolder(masterPath+"/"+entry.Name(), slavePath+"/"+entry.Name())
			}

		} else {

			wgCSF.Add(1)

			go func() {

				defer wgCSF.Done()
				_ = deleteFile(entry, masterPath, slavePath)

			}()
		}

	}
	wgCSF.Wait()
	return nil
}

// check if subfolder  exists in the source folder. If not - delete the subfolder in slave
func removeFolder(name, masterPath, slavePath string) (bool, error) {

	var logInfo logger.LogMessage = logger.LogMessage{LogType: logger.LogInfo, Ref: "removeFolder", Message: ""}
	var logError logger.LogMessage = logger.LogMessage{LogType: logger.LogError, Ref: "removeFolder", Message: ""}

	folder, err := os.ReadDir(masterPath)

	if err != nil {

		logError.Message = err.Error()
		logger.LogChan <- logError

		return false, err

	}

	for _, msEntry := range folder {

		if msEntry.Name() == name && msEntry.IsDir() {

			return false, nil
		}

	}

	_ = purgeFolder(slavePath + "/" + name)

	err = os.Remove(slavePath + "/" + name)

	if err != nil {

		logError.Message = err.Error()
		logger.LogChan <- logError

		return false, err

	}

	logInfo.Message = "Folder " + slavePath + "/" + name + " deleted"
	logger.LogChan <- logInfo

	return true, nil
}

// func check if the file exists in the source folder. If not - delete the file
func deleteFile(entry os.DirEntry, masterPath, slavePath string) error {

	if entry.IsDir() {
		return errors.New("not a file")
	}

	var logInfo logger.LogMessage = logger.LogMessage{LogType: logger.LogInfo, Ref: "deleteFile", Message: ""}
	var logError logger.LogMessage = logger.LogMessage{LogType: logger.LogError, Ref: "deleteFile", Message: ""}

	folder, err := os.ReadDir(masterPath)

	if err != nil {

		logError.Message = err.Error()
		logger.LogChan <- logError

		return err

	}

	exist := false

	for _, msEntry := range folder {

		if msEntry.Name() == entry.Name() {

			exist = true
		}

	}

	if !exist {

		err = os.Remove(slavePath + "/" + entry.Name())

		if err != nil {

			logError.Message = err.Error()
			logger.LogChan <- logError
			return err

		}

		logInfo.Message = "File " + slavePath + "/" + entry.Name() + " deleted"
		logger.LogChan <- logInfo

	}
	return nil
}

// func removes all files from a folder
func purgeFolder(path string) error {

	var logInfo logger.LogMessage = logger.LogMessage{LogType: logger.LogInfo, Ref: "purgeFolder", Message: ""}
	var logError logger.LogMessage = logger.LogMessage{LogType: logger.LogError, Ref: "purgeFolder", Message: ""}

	folder, err := os.ReadDir(path)

	if err != nil {

		logError.Message = err.Error()
		logger.LogChan <- logError
		return err

	}

	for _, entry := range folder {

		err = os.Remove(path + "/" + entry.Name())

		if err != nil {

			logError.Message = err.Error()
			logger.LogChan <- logError

		}

	}

	logInfo.Message = "Folder " + path + " is purged"
	logger.LogChan <- logInfo

	return nil

}
