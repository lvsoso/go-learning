package dirlock

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDirLock_Lock(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "*.lock-dir")
	assert.Nil(t, err)
	defer func() {
		_, err := os.Stat(tmp)
		if err != nil {
			return
		}
		os.RemoveAll(tmp)
	}()

	dirLock1 := New(tmp)
	err = dirLock1.Lock()
	if err != nil {
		t.Error(err)
	}

	go func() {
		dirLock2 := New(tmp)
		err := dirLock2.Lock()
		if err != nil {
			t.Log(err)
		}
		runtime.Gosched()
		time.Sleep(3 * time.Second)
		err = dirLock2.Lock()
		if err != nil {
			t.Error(err)
		}
		err = dirLock2.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()

	err = dirLock1.Unlock()
	if err != nil {
		t.Error(err)
	}
	time.Sleep(3 * time.Second)
}
