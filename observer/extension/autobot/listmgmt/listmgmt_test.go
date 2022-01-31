package listmgmt

import (
	"testing"

	"fort.plus/im"
	"fort.plus/repository"
)

func TestCheckResponse(t *testing.T) {
	Register(nil, "", "testlist")
	var message im.Message
	message.Text = "/testlist ls"
	if !repository.IsCallable(message) {
		t.Error("Expected response to be true for /testlist ls")
	}
}

func TestCheckResponseNotCallable(t *testing.T) {
	Register(nil, "", "testlist")
	var message im.Message
	message.Text = "/other ls"
	if repository.IsCallable(message) {
		t.Error("Expected response to be false")
	}
}
