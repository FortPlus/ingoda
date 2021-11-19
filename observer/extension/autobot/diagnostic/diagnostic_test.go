package diagnostic

import (
	"testing"

	"fort.plus/repository"
    "fort.plus/messengers/telegram"
)

func TestCheckWithIPv4Host(t *testing.T) {
    var message telegram.Message
    message.Text = "/check 127.0.0.1"
    if !repository.IsCallable(message){
        t.Error("Expected response to be true for check request with host name or address")
    }
}

func TestCheckWithIPv6Host(t *testing.T) {
    var message telegram.Message
    message.Text = "/check ::1/128"
    if !repository.IsCallable(message){
        t.Error("Expected response to be true for check request with host name or address")
    }
}

func TestCheckWithHostName(t *testing.T) {
    var message telegram.Message
    message.Text = "/check some-host-name"
    if !repository.IsCallable(message){
        t.Error("Expected response to be true for check request with host name or address")
    }
}

func TestCheckWithHostFQDN(t *testing.T) {
    var message telegram.Message
    message.Text = "/check host.example.com"
    if !repository.IsCallable(message){
        t.Error("Expected response to be true for check request with host name or address")
    }
}

func TestCheckWithoutHost(t *testing.T) {
    var message telegram.Message
    message.Text = "/check  "
    if repository.IsCallable(message){
        t.Error("Expected response to be false for check request without host name or address")
    }
}


