package helper

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type (
	IPAFactory struct {}
)

func NewIPAFactory() *IPAFactory {
	return &IPAFactory{}
}

func (i *IPAFactory) getDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("Error remove dir, path: %s, error: %s", path, err.Error())
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error create dir, path: %s, error: %s", path, err.Error())
	}

	return nil
}

//If file path empty, then create uniq temp fileName in is temp
func (i *IPAFactory) writeToFile(filePath string, file io.Reader) (string, error) {
	var fileDestPath string
	if len(filePath) == 0 {
		tmpFile, _ := ioutil.TempFile("", "goitunes_files_")
		fileDestPath = tmpFile.Name()
	} else {
		fileDestPath = filePath
	}

	filenameSource, err := os.Create(fileDestPath)
	if err != nil {
		return "", fmt.Errorf("Error create temp file, filename: %s, error: %s", fileDestPath, err.Error())
	}
	defer filenameSource.Close()

	_, err = io.Copy(filenameSource, file)
	if err != nil {
		return "", fmt.Errorf("Error copy file from reader, filename: %s, error: %s", fileDestPath, err.Error())
	}

	return fileDestPath, nil
}

func (i *IPAFactory) unzipToTmpDir(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("Source file is not exists, filepath: %s, error: %s", filePath, err.Error())
	}

	unzipDirPath, _ := ioutil.TempDir("", "goitunes_unzip_")

	cmd := exec.Command("unzip", filePath, "-d", unzipDirPath)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error unzip file, error: %s", err.Error())
	}

	return unzipDirPath, nil
}

func (i *IPAFactory) zipPath(filePath string, sourcePath string) (string, error) {
	var zipFilePath string
	if len(filePath) == 0 {
		tmpFile, _ := ioutil.TempFile("", "goitunes_zip_")
		zipFilePath = tmpFile.Name() + ".zip"
	} else {
		zipFilePath = filePath
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("`cd %s && /usr/bin/zip -qr %s %s`", sourcePath, zipFilePath, "*"))
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error zip directory, error: %s", err.Error())
	}

	return zipFilePath, nil
}

func (i *IPAFactory) removeAll(elements []string) {
	for _, elem := range elements {
		fi, err := os.Stat(elem)
		if err != nil {
			continue
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			err = os.RemoveAll(elem)
			if err != nil {
				continue
			}
		case mode.IsRegular():
			err = os.Remove(elem)
			if err != nil {
				continue
			}
		}
	}
}
