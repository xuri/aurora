package beanstalk

import (
	"testing"
	"time"
)

func TestTubeSetReserve(t *testing.T) {
	c := NewConn(mock("reserve-with-timeout 1\r\n", "RESERVED 1 1\r\nx\r\n"))
	id, body, err := c.Reserve(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatal("expected 1, got", id)
	}
	if len(body) != 1 || body[0] != 'x' {
		t.Fatalf("bad body, expected %#v, got %#v", "x", string(body))
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTubeSetReserveTimeout(t *testing.T) {
	c := NewConn(mock("reserve-with-timeout 1\r\n", "TIMED_OUT\r\n"))
	_, _, err := c.Reserve(time.Second)
	if cerr, ok := err.(ConnError); !ok {
		t.Log(err)
		t.Logf("%#v", err)
		t.Fatal("expected ConnError")
	} else if cerr.Err != ErrTimeout {
		t.Log(err)
		t.Logf("%#v", err)
		t.Fatal("expected ErrTimeout")
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}
