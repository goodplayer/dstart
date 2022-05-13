package dstartcore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

var ErrFileNotExist = errors.New("file not exist")

func ProcessDstartFromConfigFile(configFilePath string) {
	log.Println("start from config file.")

	wd, err := os.Getwd()
	if err != nil {
		panic(errors.New(fmt.Sprint("get working directory error:", err)))
	}
	wdConfig := fmt.Sprintf("%s%c%s", wd, os.PathSeparator, configFilePath)
	log.Println("guessing config file path - working directory config:", wdConfig)

	exeFile, err := os.Executable()
	if err != nil {
		panic(errors.New(fmt.Sprint("get os.Executable() error:", err)))
	}
	exeFolder := filepath.Dir(exeFile)
	exeFolderConfig := fmt.Sprintf("%s%c%s", exeFolder, os.PathSeparator, configFilePath)
	log.Println("guessing config file path - application file directory config:", exeFolderConfig)

	var fileData []byte
	if len(fileData) <= 0 {
		if data, err := readFileData(wdConfig); err == nil {
			fileData = data
			log.Println("load config file:", wdConfig)
		} else if err != ErrFileNotExist {
			panic(errors.New(fmt.Sprintf("try to read file=[%s] error:%s", wdConfig, err)))
		}
	}
	if len(fileData) <= 0 {
		if data, err := readFileData(exeFolderConfig); err == nil {
			fileData = data
			log.Println("load config file:", exeFolderConfig)
		} else if err != ErrFileNotExist {
			panic(errors.New(fmt.Sprintf("try to read file=[%s] error:%s", exeFolderConfig, err)))
		}
	}
	if len(fileData) <= 0 {
		panic(errors.New("no config file found"))
	}

	c := new(Config)
	err = toml.Unmarshal(fileData, c)
	if err != nil {
		log.Println("unmarshal config file error.", err)
		panic(errors.New("unmarshal config file error"))
	}

	err = startFromConfig(c)
	if err != nil {
		log.Println("start from config error.", err)
		panic(errors.New("start from config error"))
	}
}

func readFileData(file string) ([]byte, error) {
	var fileData []byte

	fs := afero.NewOsFs()
	if ok, err := afero.Exists(fs, file); err != nil {
		panic(errors.New(fmt.Sprintf("check file=[%s] exist error:%s", file, err)))
	} else if ok {
		if file, err := fs.Open(file); err != nil {
			panic(errors.New(fmt.Sprintf("open file=[%s] error:%s", file, err)))
		} else {
			func() {
				defer file.Close()
				fileData, err = ioutil.ReadAll(file)
				if err != nil {
					panic(errors.New(fmt.Sprintf("readAll file=[%s] error:%s", file, err)))
				}
			}()
		}
	} else {
		return nil, ErrFileNotExist
	}

	return fileData, nil
}
