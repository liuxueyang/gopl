package main

import (
	"fmt"
	"io"
	"os"
)

type ReadString struct {
	s string
	i int64
}

func (r *ReadString) Read(p []byte) (n int, err error) {
	n = copy(p, r.s[r.i:])
	r.i += int64(n)
	if r.i >= int64(len(r.s)) {
		err = io.EOF
	}
	return
}

func NewReader(s string) io.Reader {
	return &ReadString{s: s, i: 0}
}

type LimitRead struct {
	r io.Reader
	n int64 // read bytes limit
	i int64 // current read position
}

func (l *LimitRead) Read(p []byte) (n int, err error) {
	if l.i >= l.n {
		return 0, io.EOF
	}

	remain := l.n - l.i
	if remain > int64(len(p)) {
		remain = int64(len(p))
	}
	n, err = l.r.Read(p[:remain])

	l.i += int64(n)
	if l.i == l.n {
		err = io.EOF
	}
	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitRead{r, n, 0}
}

func main() {

	reader := NewReader("abcde")
	limitReader := LimitReader(reader, 5)

	p := make([]byte, 3, 10)
	n, err := limitReader.Read(p)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(n, err, string(p[:n]))

	p1 := make([]byte, 10, 10)
	n, err = limitReader.Read(p1)
	fmt.Println(n, err, string(p1[:n]))
}
