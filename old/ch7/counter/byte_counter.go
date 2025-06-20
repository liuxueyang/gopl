package counter

import (
	"strconv"
)

type ByteCounter int

func (b *ByteCounter) Write(p []byte) (int, error) {
	*b += ByteCounter(len(p))
	return len(p), nil
}

func (b ByteCounter) String() string {
	return strconv.FormatInt(int64(b), 10)
}
