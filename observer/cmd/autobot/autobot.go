package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"fort.plus/extension/autobot/chuck"
	"fort.plus/extension/autobot/listmgmt"
	"fort.plus/fperror"
	"fort.plus/im"
	"fort.plus/im/telegram"
	"fort.plus/repository"
)

func main() {
	log.SetOutput(os.Stdout)

	app := kingpin.New(filepath.Base(os.Args[0]), "telegram bot")
	token := app.Flag("telegram.token", "Telegram token.").Required().OverrideDefaultFromEnvar("TELEGRAM_TOKEN").String()
	listServerUri := app.Flag("list.manager.uri", "URI of the server keeping lists.").Default("http://localhost:9190").OverrideDefaultFromEnvar("LIST_MANAGER_URI").String()

	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var t im.Carrier = telegram.New(*token)


	chuck.Register(t)
	listmgmt.Register(t, *listServerUri, "ban")

	for {
		messages, err := t.GetMessages()
		if err != nil {
			log.Println(fperror.Warning("Can't get message from telegram", nil))
		}
		for _, message := range messages {
			go repository.Call(message)
		}
		time.Sleep(20 * time.Second)
	}
}
