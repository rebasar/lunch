package fsutil

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

func GetAbsoluteFilename(filename string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(currentUser.HomeDir, filename), nil
}

func ReadConfigFile(fullPath string) ([]byte, error) {
	fp, err := os.Open(fullPath)
	defer fp.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(fp)
}

func WriteConfigFile(body []byte, fullPath string) error {
	return ioutil.WriteFile(fullPath, body, 0666)
}

func RemoveConfigFile(fullPath string) error {
	return os.Remove(fullPath)
}
