package storage

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocalStorage(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "*.lock-dir")
	assert.Nil(t, err)
	defer func() {
		_, err := os.Stat(tmp)
		if err != nil {
			return
		}
		os.RemoveAll(tmp)
	}()

	f1, err := ioutil.TempFile(tmp, "*.local-storage-test")
	assert.Nil(t, err)
	dst1 := f1.Name() + ".dst"

	f2, err := ioutil.TempFile(tmp, "*.local-storage-test")
	assert.Nil(t, err)
	dst2 := f2.Name() + ".dst"

	var ls Storage = NewLocalStorage(tmp)
	ls.Move(f1.Name(), dst1)

	go func() {
		ls.Move(f2.Name(), dst2)
	}()

	time.Sleep(1 * time.Second)

	exist, err := ls.Exist(dst1)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = ls.Exist(dst2)
	assert.Nil(t, err)
	assert.True(t, exist)
}
