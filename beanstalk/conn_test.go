package beanstalk

import (
	"testing"
	"time"
)

func TestNameTooLong(t *testing.T) {
	c := NewConn(mock("", ""))

	tube := NewTube(c, string(make([]byte, 201)))
	_, err := tube.Put([]byte("foo"), 0, 0, 0)
	if e, ok := err.(NameError); !ok || e.Err != ErrTooLong {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNameEmpty(t *testing.T) {
	c := NewConn(mock("", ""))

	tube := NewTube(c, "")
	_, err := tube.Put([]byte("foo"), 0, 0, 0)
	if e, ok := err.(NameError); !ok || e.Err != ErrEmpty {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNameBadChar(t *testing.T) {
	c := NewConn(mock("", ""))

	tube := NewTube(c, "*")
	_, err := tube.Put([]byte("foo"), 0, 0, 0)
	if e, ok := err.(NameError); !ok || e.Err != ErrBadChar {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNegativeDuration(t *testing.T) {
	c := NewConn(mock("", ""))
	tube := NewTube(c, "foo")
	for _, d := range []time.Duration{-100 * time.Millisecond, -2 * time.Second} {
		if _, err := tube.Put([]byte("hello"), 0, d, d); err == nil {
			t.Fatalf("put job with negative duration %v expected error, got nil", d)
		}
	}
}

func TestDeleteMissing(t *testing.T) {
	c := NewConn(mock("delete 1\r\n", "NOT_FOUND\r\n"))

	err := c.Delete(1)
	if e, ok := err.(ConnError); !ok || e.Err != ErrNotFound {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestUse(t *testing.T) {
	c := NewConn(mock(
		"use foo\r\nput 0 0 0 5\r\nhello\r\n",
		"USING foo\r\nINSERTED 1\r\n",
	))
	tube := NewTube(c, "foo")
	id, err := tube.Put([]byte("hello"), 0, 0, 0)
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

func TestWatchIgnore(t *testing.T) {
	c := NewConn(mock(
		"watch foo\r\nignore default\r\nreserve-with-timeout 1\r\n",
		"WATCHING 2\r\nWATCHING 1\r\nRESERVED 1 1\r\nx\r\n",
	))
	ts := NewTubeSet(c, "foo")
	id, body, err := ts.Reserve(time.Second)
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

func TestBury(t *testing.T) {
	c := NewConn(mock("bury 1 3\r\n", "BURIED\r\n"))

	err := c.Bury(1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTubeKickJob(t *testing.T) {
	c := NewConn(mock("kick-job 3\r\n", "KICKED\r\n"))

	err := c.KickJob(3)
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	c := NewConn(mock("delete 1\r\n", "DELETED\r\n"))

	err := c.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestListTubes(t *testing.T) {
	c := NewConn(mock("list-tubes\r\n", "OK 14\r\n---\n- default\n\r\n"))

	l, err := c.ListTubes()
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 1 || l[0] != "default" {
		t.Fatalf("expected %#v, got %#v", []string{"default"}, l)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestPeek(t *testing.T) {
	c := NewConn(mock("peek 1\r\n", "FOUND 1 1\r\nx\r\n"))

	body, err := c.Peek(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body[0] != 'x' {
		t.Fatalf("bad body, expected %#v, got %#v", "x", string(body))
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestPeekTwice(t *testing.T) {
	c := NewConn(mock(
		"peek 1\r\npeek 1\r\n",
		"FOUND 1 1\r\nx\r\nFOUND 1 1\r\nx\r\n",
	))

	body, err := c.Peek(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body[0] != 'x' {
		t.Fatalf("bad body, expected %#v, got %#v", "x", string(body))
	}

	body, err = c.Peek(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body[0] != 'x' {
		t.Fatalf("bad body, expected %#v, got %#v", "x", string(body))
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRelease(t *testing.T) {
	c := NewConn(mock("release 1 3 2\r\n", "RELEASED\r\n"))

	err := c.Release(1, 3, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestStats(t *testing.T) {
	c := NewConn(mock("stats\r\n", "OK 10\r\n---\na: ok\n\r\n"))

	m, err := c.Stats()
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

func TestStatsJob(t *testing.T) {
	c := NewConn(mock("stats-job 1\r\n", "OK 10\r\n---\na: ok\n\r\n"))

	m, err := c.StatsJob(1)
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

func TestTouch(t *testing.T) {
	c := NewConn(mock("touch 1\r\n", "TOUCHED\r\n"))

	err := c.Touch(1)
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}
