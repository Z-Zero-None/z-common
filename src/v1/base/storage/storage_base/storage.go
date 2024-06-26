package storage_base

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"z-common/src/v1/base/storage"
	"z-common/src/v1/base/utils"
)

type baseStorage struct {
	baseRoot  string
	storePath string
}

func NewBaseStorage() storage.IStorage {
	return &baseStorage{
		baseRoot:  os.Getenv("base_root"),
		storePath: os.Getenv("store_path"),
	}
}

func (s *baseStorage) checkNoExistFilePath(filename string) (string, error) {
	filePath := filepath.Join(s.baseRoot, s.storePath, filename)
	fileDir := filepath.Dir(filePath)
	if !utils.DirExist(fileDir) {
		err := s.createDir(fileDir)
		if err != nil {
			errMsg := fmt.Sprintf("Create dir for file |%s|  failed.", filePath)
			return "", errors.Wrap(err, errMsg)
		}
	}
	return filePath, nil
}
func (s *baseStorage) createDir(dir string) error {
	if len(dir) == 0 {
		return errors.New("Dir is empty.")
	}
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Make dir |%s| failed.", dir))
	}
	return nil
}

func (s *baseStorage) Upload(filename string, reader io.Reader) (string, error) {
	// get filePath: baseRoot+storePath+filename
	savePath, err := s.checkNoExistFilePath(filename)
	if err != nil {
		return "", errors.Wrap(err, "Check filepath failed.")
	}
	targetFile, err := os.Create(savePath)
	if err != nil {
		errMsg := fmt.Sprintf("Open target file |%s| failed.", savePath)
		return "", errors.Wrap(err, errMsg)
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, reader)
	if err != nil {
		errMsg := fmt.Sprintf("Save file |%s|  failed.", savePath)
		return "", errors.Wrap(err, errMsg)
	}
	return savePath, nil
}

func (s *baseStorage) checkExistFilePath(filePath string) error {
	fileDir := filepath.Dir(filePath)
	prefixDir := filepath.Join(s.baseRoot, s.storePath)
	if strings.HasPrefix(fileDir, prefixDir) {
		return nil /**/
	}
	errMsg := fmt.Sprintf("File |%s|  isn't in correct dir |%s|", filePath, prefixDir)
	return errors.New(errMsg)
}

func (s *baseStorage) Load(filePath string) (io.ReadCloser, error) {
	err := s.checkExistFilePath(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Check filepath failed.")
	}
	file, err := os.Open(filePath)
	if err != nil {
		errMsg := fmt.Sprintf("Read file |%s| failed when loading file reader.", filePath)
		return nil, errors.WithMessage(err, errMsg)
	}
	return file, nil
}
func (s *baseStorage) Remove(filePath string) error {
	err := s.checkExistFilePath(filePath)
	if err != nil {
		return errors.Wrap(err, "Check filepath failed.")
	}
	err = os.Remove(filePath)
	if err != nil {
		errMsg := fmt.Sprintf("Delete file |%s| failed when deleting.", filePath)
		return errors.WithMessage(err, errMsg)
	}
	return nil
}

func (s *baseStorage) Storage() {

}
