package mbroker

import (
	"testing"
	"log"
)

func TestWaitForSubject(t *testing.T) {
    subject := "someid."+GetBatchId("somestring")
    log.Printf("wait for subject:%s", subject)
    WaitForSubject(subject, 10)
}

func TestSendMessage(t *testing.T) {
    subject := "someid."+GetBatchId("somestring")
    log.Printf("send message with subject:%s", subject)
    SendMessage(subject, "some data")
}