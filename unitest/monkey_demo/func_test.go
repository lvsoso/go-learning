package monkey_demo

import (
	"strings"
	"testing"
	"varys"

	"bou.ke/monkey"
)

// go test -run=TestMyFunc -v -gcflags=-l
func TestMyFunc(t *testing.T) {
	monkey.Patch(varys.GetInfoByUID, func(int64) (*varys.UserInfo, error) {
		return &varys.UserInfo{Name: "lvsoso"}, nil
	})

	ret := MyFunc(123)
	if !strings.Contains(ret, "lvsoso") {
		t.Fatal()
	}
}
