package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNameFromPath(t *testing.T) {
	assert.Equal(t, "abc", GetNameFromPath("abc/"))
	assert.Equal(t, "abc", GetNameFromPath("abc/"))
	assert.Equal(t, "def/", GetNameFromPath("abc/def/"))
}

func TestEncodePath(t *testing.T) {
	expected := "c5451a0a-b500-4162-8eb1-c87651d01440%20!%40%23%24%25%5E%26*()_%2B-%3D%3B'%3E%3C%2C%3F.txt"
	assert.Equal(t, expected, EncodePath("c5451a0a-b500-4162-8eb1-c87651d01440 !@#$%^&*()_+-=;'><,?.txt"))

	expected = "%3B%2C/%3F%3A%40%26%3D%2B%24"
	assert.Equal(t, expected, EncodePath(";,/?:@&=+$"))

	assert.Equal(t, "-_.!~*'()", EncodePath("-_.!~*'()"))

	assert.Equal(t, "%23", EncodePath("#"))

	assert.Equal(t, "ABC%20abc%20123", EncodePath("ABC abc 123"))

	assert.Equal(t, "%E4%BD%A0%E5%A5%BD%EF%BC%8C%E4%B8%96%E7%95%8C%EF%BC%81%E2%9D%A4", EncodePath("你好，世界！❤"))
}
