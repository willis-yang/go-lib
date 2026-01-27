package intervalScheduler

import (
	"fmt"
	"testing"
)

func TestMustNewScheduler(t *testing.T) {
	got, err := MustNewScheduler(10)
	if err != nil {
		t.Errorf("MustNewScheduler() error = %v", err)
		return
	}
	got.AddFunc(func() {
		fmt.Println("testing interval scheduler")
	})
	got.Start()
}
