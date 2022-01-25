package main

import (
	"log"
	"os"
	"time"

	_ "fort.plus/extension/autobot"
	"fort.plus/fperror"
	"fort.plus/messengers/telegram"
	"fort.plus/repository"
)

func main() {
	log.SetOutput(os.Stdout)
	for {
		messages, err := telegram.GetMessages()
		if err != nil {
			log.Println(fperror.Warning("Can't get message from telegram", nil))
		}
		for _, message := range messages {
			go repository.Call(message)
		}
		time.Sleep(20 * time.Second)
	}
}
