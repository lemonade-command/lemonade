package lemon

import "testing"

func TestConvertLineEnding(t *testing.T) {
	assert := func(text, option, expected string) {
		if got := ConvertLineEnding(text, option); got != expected {
			t.Errorf("Expected: %+v, got %+v", []byte(expected), []byte(got))
		}
	}

	assert("aaa\r\nbbb", "lf", "aaa\nbbb")
	assert("aaa\rbbb", "lf", "aaa\nbbb")
	assert("aaa\nbbb", "lf", "aaa\nbbb")

	assert("aaa\r\nbbb", "crlf", "aaa\r\nbbb")
	assert("aaa\rbbb", "crlf", "aaa\r\nbbb")
	assert("aaa\nbbb", "crlf", "aaa\r\nbbb")

	assert("a\r", "crlf", "a\r\n")
	assert("\na", "crlf", "\r\na")
}
