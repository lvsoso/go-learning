package monkey_demo

import (
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"
)

func TestUser_GetInfo(t *testing.T) {
	var u = &User{
		Name:     "lvsoso",
		Birthday: "2000-01-01",
	}

	// 为对象方法打桩
	monkey.PatchInstanceMethod(reflect.TypeOf(u), "CalcAge", func(*User) int {
		return 18
	})

	ret := u.GetInfo()
	if strings.Contains(ret, "!") {
		t.Fatal()
	}
}
