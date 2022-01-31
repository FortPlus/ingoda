package diff

import (
	"runtime"
	"testing"
)

func trace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

func TestSameStringsDistance(t *testing.T) {
	Reset()
	GetLevenshteinDistance("abc")
	b := GetLevenshteinDistance("abc")
	if b != 0 {
		t.Errorf("%s, expected 0 distance for same string, got%d", trace(), b)
	}
}

func TestDifferentStringsDistance(t *testing.T) {
	Reset()
	GetLevenshteinDistance("abc")
	b := GetLevenshteinDistance("abx")
	if b != 1 {
		t.Errorf("%s, expected 1, got:%d", trace(), b)
	}
}

// If set is empty, then new message will exceed threshold
func TestEmptySamplesThreshold(t *testing.T) {
	Reset()
	if !IsThresholdExceeded("first time message", 10) {
		t.Errorf("%s, threshold exceeded with empty samples set", trace())
	}
}

func TestThresholdNotExceeded(t *testing.T) {
	Reset()
	GetLevenshteinDistance("abc")
	if IsThresholdExceeded("abcde", 2) {
		t.Errorf("%s, unexpected exceeded threshold", trace())
	}
}

func TestThresholdExceeded(t *testing.T) {
	Reset()
	GetLevenshteinDistance("abc")
	if !IsThresholdExceeded("abcde", 20) {
		t.Errorf("%s, unexpected exceeded threshold", trace())
	}
}
