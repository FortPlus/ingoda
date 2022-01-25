package banmgmt

import (
	"testing"

	"fort.plus/messengers/telegram"
	"fort.plus/repository"
)

func TestCheckWithoutHost(t *testing.T) {
	var message telegram.Message
	message.Text = "/ban ls"
	if !repository.IsCallable(message) {
		t.Error("Expected response to be true for /ban ls")
	}
}
