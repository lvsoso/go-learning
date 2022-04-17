package monkey_demo

import (
	"fmt"
	"varys"
)

func MyFunc(uid int64) string {
	u, err := varys.GetInfoByUID(uid)
	if err != nil {
		return "welcome"
	}

	fmt.Printf("-->%#v\n", u)
	return fmt.Sprintf("hello %s\n", u.Name)
}
