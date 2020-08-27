package counting_writer

import (
	"fmt"
	"gitee.com/liuxueyang/gopl/ch7/counter"
	"testing"
)

func TestCountingWriter(t *testing.T) {
	var w counter.WordCounter
	w1, b := CountingWriter(&w)

	fmt.Fprintf(w1, "a bb")
	t.Logf("w=%s, b=%d", &w, *b)

	fmt.Fprintf(w1, "ccd")
	t.Logf("w=%s, b=%d", &w, *b)
}
