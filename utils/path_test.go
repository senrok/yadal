package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNameFromPath(t *testing.T) {
	assert.Equal(t, "abc", GetNameFromPath("abc/"))
	assert.Equal(t, "abc", GetNameFromPath("abc/"))
	assert.Equal(t, "def", GetNameFromPath("abc/def/"))
}
