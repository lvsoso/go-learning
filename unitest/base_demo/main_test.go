package base_demo

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("write setup code here ...")
	// flags.Parse()
	retCode := m.Run()
	fmt.Println("write teardown code here ...")
	os.Exit(retCode)
}
