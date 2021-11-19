package config

import (
	"testing"
)

func TestGetCurrentImmutable(t *testing.T) {
    configData1 := GetCurrent()
    configData2 := GetCurrent()

    configData1.TelegramToken = "1"
    configData2.TelegramToken = "2"

    if configData1.TelegramToken == configData2.TelegramToken {
        t.Error("config data is not immutable")
    }
}
