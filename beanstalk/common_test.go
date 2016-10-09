package beanstalk

import (
	"fmt"
	"io"
	"strings"
)

type mockError struct {
	exp []byte
	got []byte
}

func (e mockError) Error() string {
	return fmt.Sprintf(
		"mock error: exp %#v, got %#v",
		string(e.exp),
		string(e.got),
	)
}

type mockIO struct {
	recv *strings.Reader
	send *strings.Reader
}

func mock(recv, send string) io.ReadWriteCloser {
	return &mockIO{strings.NewReader(recv), strings.NewReader(send)}
}

func (m mockIO) Read(b []byte) (int, error) {
	return m.send.Read(b)
}

func (m mockIO) Write(got []byte) (n int, err error) {
	exp := make([]byte, len(got))
	n, err = m.recv.Read(exp)
	if err != nil {
		return n, err
	}
	exp = exp[:n]
	for i := range exp {
		if exp[i] != got[i] {
			return i, mockError{exp, got}
		}
	}
	if n != len(got) {
		return n, mockError{exp, got}
	}
	return n, err
}

func (m mockIO) Close() error {
	if m.recv.Len() == 0 && m.send.Len() == 0 {
		return nil
	}
	if m.recv.Len() > 0 {
		b := make([]byte, m.recv.Len())
		m.recv.Read(b)
		return mockError{b, nil}
	}
	b := make([]byte, m.send.Len())
	m.send.Read(b)
	return mockError{b, nil}
}
