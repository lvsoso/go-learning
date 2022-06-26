package storage

import (
	"local-storage/dirlock"
	"os"
	"path/filepath"
	"time"
)

type Storage interface {
	Move(src string, dst string) error
	Exist(dst string) (bool, error)
}

type LocalStorage struct {
	dirPath string
}

func NewLocalStorage(dirPath string) *LocalStorage {
	return &LocalStorage{
		dirPath: dirPath,
	}
}

func (ls *LocalStorage) Move(src string, dst string) error {
	dl := dirlock.New(ls.dirPath)
	wait := 10
	err := dl.Lock()
	for err != nil {
		time.Sleep(1 * time.Second)
		wait--
		if wait == 0 {
			return err
		}
		continue
	}
	defer dl.Unlock()
	newDst := filepath.Join(ls.dirPath, dst)
	return os.Rename(src, newDst)
}

func (ls *LocalStorage) Exist(dst string) (bool, error) {
	newDst := filepath.Join(ls.dirPath, dst)
	_, err := os.Stat(newDst)
	if err != nil {
		return false, err
	}
	return true, nil
}
