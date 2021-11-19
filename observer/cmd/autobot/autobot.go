package main

import (
	"log"
	"time"
    "os"
	_ "fort.plus/extension/autobot"
	"fort.plus/messengers/telegram"
	"fort.plus/repository"
	"fort.plus/fperror"
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
		log.Println("got message:",messages)
		time.Sleep(20 * time.Second)
	}
}
