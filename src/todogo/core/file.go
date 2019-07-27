package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadBytes(fpath string) ([]byte, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var bytes []byte
	bytes, err = ioutil.ReadAll(file)
	return bytes, err
}

func WriteBytes(fpath string, bytes []byte) error {
	err := CheckAndMakeDir(filepath.Dir(fpath))
	if err != nil {
		return err
	}
	file, err := os.Create(fpath)
	defer file.Close()
	if err != nil {
		return err
	}
	file.Write(bytes)
	file.Sync()
	return nil
}

// PathExists returns true if the path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, err
}

// CheckAndMakeDir check if the directory exists and creates it with all its
// tree structure if not.
func CheckAndMakeDir(dpath string) error {
	exists, err := PathExists(dpath)
	if exists && err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(dpath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
