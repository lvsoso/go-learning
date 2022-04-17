package monkey_demo

import (
	"fmt"
	"time"
)

type User struct {
	Name     string
	Birthday string
}

// CalcAge 计算用户年龄
func (u *User) CalcAge() int {
	t, err := time.Parse("2006-01-02", u.Birthday)
	if err != nil {
		return -1
	}
	return int(time.Now().Sub(t).Hours()/24.0) / 365
}

// GetInfo 获取用户相关信息
func (u *User) GetInfo() string {
	age := u.CalcAge()
	if age <= 0 {
		return fmt.Sprintf("%s %d 岁 ！！！。", u.Name, age)
	}
	return fmt.Sprintf("%s今年%d岁了。", u.Name, age)
}
