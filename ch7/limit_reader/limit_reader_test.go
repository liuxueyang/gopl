package main

import (
	"io"
	"testing"
)

func TestLimitReader(t *testing.T) {

	reader := NewReader("abcde")
	limitReader := LimitReader(reader, 5)

	p := make([]byte, 3, 10)
	n, err := limitReader.Read(p)

	if n != 3 || err != nil || string(p) != "abc" {
		t.Errorf("n=%v, err=%v", n, err)
		return
	}
}

func TestLimitReader2(t *testing.T) {

	reader := NewReader("abcde")
	limitReader := LimitReader(reader, 30)

	p := make([]byte, 3, 10)
	n, err := limitReader.Read(p)

	if n != 3 || err != nil || string(p) != "abc" {
		t.Errorf("n=%v, err=%v", n, err)
		return
	}
}

func TestLimitReader3(t *testing.T) {

	reader := NewReader("abcde")
	limitReader := LimitReader(reader, 30)

	p := make([]byte, 3, 10)
	n, err := limitReader.Read(p)

	if n != 3 || err != nil || string(p) != "abc" {
		t.Errorf("n=%v, err=%v", n, err)
		return
	}

	p1 := make([]byte, 10, 10)
	n, err = limitReader.Read(p1)
	p1s := string(p1[:n])

	if err != io.EOF || n != 2 || string(p1[:n]) != "de" {
		t.Errorf("n=%v, err=%v, p1s=%v", n, err, p1s)
		return
	}
}
