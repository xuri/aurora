package beanstalk

import (
	"testing"
	"time"
)

func TestTubePut(t *testing.T) {
	c := NewConn(mock("put 0 0 0 3\r\nfoo\r\n", "INSERTED 1\r\n"))

	id, err := c.Put([]byte("foo"), 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatal("expected 1, got", id)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTubePeekReady(t *testing.T) {
	c := NewConn(mock("peek-ready\r\n", "FOUND 1 1\r\nx\r\n"))

	id, body, err := c.PeekReady()
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

func TestTubePeekDelayed(t *testing.T) {
	c := NewConn(mock("peek-delayed\r\n", "FOUND 1 1\r\nx\r\n"))

	id, body, err := c.PeekDelayed()
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

func TestTubePeekBuried(t *testing.T) {
	c := NewConn(mock("peek-buried\r\n", "FOUND 1 1\r\nx\r\n"))

	id, body, err := c.PeekBuried()
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

func TestTubeKick(t *testing.T) {
	c := NewConn(mock("kick 2\r\n", "KICKED 1\r\n"))

	n, err := c.Kick(2)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("expected 1, got", n)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTubeStats(t *testing.T) {
	c := NewConn(mock("stats-tube default\r\n", "OK 10\r\n---\na: ok\n\r\n"))

	m, err := c.Tube.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 1 || m["a"] != "ok" {
		t.Fatalf("expected %#v, got %#v", map[string]string{"a": "ok"}, m)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTubePause(t *testing.T) {
	c := NewConn(mock("pause-tube default 5\r\n", "PAUSED\r\n"))

	err := c.Pause(5 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}
