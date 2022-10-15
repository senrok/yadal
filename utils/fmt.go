package utils

import "fmt"

func FmtPtrUint64(n *uint64) string {
	if n != nil {
		return fmt.Sprintf("%d", *n)
	}
	return "nil"
}
