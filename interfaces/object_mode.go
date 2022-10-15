package interfaces

import "strings"

type ObjectMode int8

const (
	Unknown ObjectMode = iota
	FILE
	DIR
)

func ObjectModeFromPath(path string) ObjectMode {
	if strings.HasSuffix(path, "/") {
		return DIR
	}
	return FILE
}

var mode2Str = []string{"Unknown", "file", "dir"}

func (m ObjectMode) IsFile() bool {
	return m == FILE
}

func (m ObjectMode) IsDir() bool {
	return m == DIR
}

func (m ObjectMode) String() string {
	return mode2Str[m]
}
