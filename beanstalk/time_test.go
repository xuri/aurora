package beanstalk

import (
	"fmt"
	"testing"
)

func TestFormatDuration(t *testing.T) {
	var d dur = 100e9
	s := fmt.Sprint(d)
	if s != "100" {
		t.Fatal("got", s, "expected 100")
	}
}
