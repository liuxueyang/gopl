package counting_writer

import (
	"io"
)

type CountWriter struct {
	a io.Writer
	c *int64
}

func (c *CountWriter) Write(p []byte) (int, error) {
	n, err := c.a.Write(p)
	*c.c += int64(n)

	return len(p), err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var c int64
	cw := CountWriter{w, &c}
	return &cw, cw.c
}
