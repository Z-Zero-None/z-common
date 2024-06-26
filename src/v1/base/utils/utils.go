package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

func DirExist(dirName string) bool {
	fileInfo, err := os.Stat(dirName)
	return (err == nil || os.IsExist(err)) && fileInfo.IsDir()
}

func RemoveDir(dirPath string) error {
	if dirPath == "" || dirPath == "/" {
		return errors.New(fmt.Sprintf("Invalid remove dir: |%s|.", dirPath))
	}
	err := os.RemoveAll(dirPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Remove dir |%s| failed.", dirPath))
	}
	return nil
}
