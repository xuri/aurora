package beanstalk

import (
	"reflect"
	"testing"
)

func TestParseDict(t *testing.T) {
	d := parseDict([]byte("---\na: 1\nb: 2\n"))
	if !reflect.DeepEqual(d, map[string]string{"a": "1", "b": "2"}) {
		t.Fatalf("got %v", d)
	}
}

func TestParseDictEmpty(t *testing.T) {
	d := parseDict([]byte{})
	if !reflect.DeepEqual(d, map[string]string{}) {
		t.Fatalf("got %v", d)
	}
}

func TestParseDictNil(t *testing.T) {
	d := parseDict(nil)
	if d != nil {
		t.Fatalf("got %v", d)
	}
}

func TestParseList(t *testing.T) {
	l := parseList([]byte("---\n- 1\n- 2\n"))
	if !reflect.DeepEqual(l, []string{"1", "2"}) {
		t.Fatalf("got %v", l)
	}
}

func TestParseListEmpty(t *testing.T) {
	l := parseList([]byte{})
	if !reflect.DeepEqual(l, []string{}) {
		t.Fatalf("got %v", l)
	}
}

func TestParseListNil(t *testing.T) {
	l := parseList(nil)
	if l != nil {
		t.Fatalf("got %v", l)
	}
}
