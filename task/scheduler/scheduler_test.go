package scheduler

import (
	"fmt"
	"testing"
)

func TestMustNewScheduler(t *testing.T) {
	c := MustNewScheduler()
	//注册报表任务
	c.AddFunc("0/10 * * * * *", func() {
		fmt.Println("exec task")
	})
}
