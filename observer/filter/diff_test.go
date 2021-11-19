package diff

import (
	"runtime"
	"testing"
	"fort.plus/fperror"
)

func trace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	//return fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
	return frame.Function
}

func TestSameStringsDistance(t *testing.T) {
     e := fperror.Throw("some shit happen")
     e.AsString()
	Reset()
	a := GetLevenshteinDistance("abc")
	b := GetLevenshteinDistance("abc")
	if a != b {
		t.Errorf("%s, %d!=%d", trace(), a, b)
	}
}

func TestDifferentStringsDistance(t *testing.T) {
	Reset()
	_ = GetLevenshteinDistance("abc")
	b := GetLevenshteinDistance("abx")
	if b != 1 {
		t.Errorf("%s, expected 1, got:%d", trace(), b)
	}
}

func TestEmptySamplesThreshold(t *testing.T) {
	Reset()
	if IsThresholdExceeded("first time message", 0) {
		t.Errorf("%s, threshold exceeded with empty samples set", trace())
	}
}

func TestEmptySamplesThresholdNotNull(t *testing.T) {
	Reset()
	if IsThresholdExceeded("first time message", 15) {
		t.Errorf("%s, threshold exceeded with empty samples set", trace())
	}
}

func TestThresholdNotExceeded(t *testing.T) {
	Reset()
	_ = GetLevenshteinDistance("abc")
	if IsThresholdExceeded("abcde", 2) {
		t.Errorf("%s, unexpected exceeded threshold", trace())
	}
}

func TestThresholdExceeded(t *testing.T) {
	Reset()
	_ = GetLevenshteinDistance("abc")
	if !IsThresholdExceeded("abcde", 20) {
		t.Errorf("%s, unexpected exceeded threshold", trace())
	}
}
