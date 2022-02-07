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

	app := kingpin.New(filepath.Base(os.Args[0]), "telegram IM bot")
	token := app.Flag("telegram.token", "Telegram token.").Required().OverrideDefaultFromEnvar("TELEGRAM_TOKEN").String()
	listServerUri := app.Flag("list.manager.uri", "URI of the server keeping lists.").Default("http://localhost:9190").OverrideDefaultFromEnvar("LIST_MANAGER_URI").String()
	banGroup := app.Flag("ban.group", "Name of banned list , managed by IM bot.").Default("ban").OverrideDefaultFromEnvar("BAN_GROUP").String()
	notifyGroup := app.Flag("notify.group", "Name of notification list, managed by IM bot.").Default("notify").OverrideDefaultFromEnvar("NOTIFY_GROUP").String()

	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var carrier im.Carrier = telegram.New(*token)

	chuck.Register(carrier)
	listmgmt.Register(carrier, *listServerUri, *banGroup)
	listmgmt.Register(carrier, *listServerUri, *notifyGroup)

	for {
		messages, err := carrier.GetMessages()
		if err != nil {
			log.Println(fperror.Warning("Can't get message from telegram", nil))
		}
		for _, message := range messages {
			go repository.Call(message)
		}
		time.Sleep(20 * time.Second)
	}
}
