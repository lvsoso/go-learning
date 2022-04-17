package gostub_demo

import (
	"testing"

	"github.com/prashantv/gostub"
)

func TestGetConfig(t *testing.T) {
	// 替换全局变量
	stubs := gostub.Stub(&configFile, "./test.toml")
	defer stubs.Reset()

	data, err := GetConfig()
	if err != nil {
		t.Fatal()
	}

	t.Logf("data:%s\n", data)
}

func TestShowNumber(t *testing.T) {
	stubs := gostub.Stub(&maxNum, 20) // 设定全局变量
	defer stubs.Reset()
	// 下面是一些测试的代码
	res := ShowNumber()
	if res != 20 {
		t.Fatal()
	}
}
