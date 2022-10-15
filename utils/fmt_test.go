package utils

import (
	"fmt"
	"testing"
)

func TestFmtPtr(t *testing.T) {
	fmt.Printf("%v\n", FmtPtrUint64(nil))
	size := uint64(2)
	fmt.Printf("%v\n", FmtPtrUint64(&size))
}
