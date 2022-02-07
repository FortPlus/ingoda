package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"gopkg.in/alecthomas/kingpin.v2"

	"fort.plus/extension/sniffer/notify"
	"fort.plus/im/telegram"
	"fort.plus/listmanager"
	"fort.plus/repository"

	"fort.plus/mbroker"
)

const (
	NOTIFY_LIST_WAIT_TIME = 10
)

var (
	app             = kingpin.New(filepath.Base(os.Args[0]), "Nats message sniffer")
	token           = app.Flag("telegram.token", "Telegram token.").Required().OverrideDefaultFromEnvar("TELEGRAM_TOKEN").String()
	telegramGroup   = app.Flag("telegram.group", "Telegram notification group.").Required().OverrideDefaultFromEnvar("TELEGRAM_GROUP").String()
	th              = app.Flag("ld.th", "Levenshtein distance threshold.").Default("14").OverrideDefaultFromEnvar("LD_TH").Int()
	natsURI         = app.Flag("nats.uri", "URI of the NATS server.").Default("nats://localhost:4222").OverrideDefaultFromEnvar("NATS_URI").String()
	subj            = app.Flag("nats.subj", "NATS message subject").Default("slg.>").OverrideDefaultFromEnvar("NATS_SUBJ").String()
	listServerUri   = app.Flag("list.manager.uri", "URI of the server keeping lists.").Default("http://localhost:9190").OverrideDefaultFromEnvar("LIST_MANAGER_URI").String()
	groupName       = app.Flag("list.manager.group", "Name of group with banned list.").Default("ban").OverrideDefaultFromEnvar("LIST_MANAGER_GROUP").String()
	notifyGroupName = app.Flag("list.manager.notify.group", "Name of group with notification patterns.").Default("notify").OverrideDefaultFromEnvar("LIST_MANAGER_NOTIFY_GROUP").String()
	importPeriod    = app.Flag("list.manager.import.period", "Number of seconds between periodical import list from server.").Default("180").OverrideDefaultFromEnvar("LIST_MANAGER_IMPORT_PERIOD").Int()
)

var wg sync.WaitGroup

func main() {
	log.SetOutput(os.Stdout)

	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.Println("prepare telegram bot")
	carrier := telegram.New(*token)

	log.Println("load records from list manager")
	listRecords := listmanager.New(*groupName)
	go listRecords.PeriodicImportFromServer(*listServerUri+"/api/v1/", *importPeriod)

	log.Println("load notification patterns from list manager")
	listNotify := listmanager.New(*notifyGroupName)
	go listNotify.PeriodicImportFromServer(*listServerUri+"/api/v1/", *importPeriod)

	for listNotify.IsEmpty() {
		log.Println("notify list is empty, wait ", NOTIFY_LIST_WAIT_TIME, "seconds and try again")
		time.Sleep(NOTIFY_LIST_WAIT_TIME * time.Second)
	}

	notify.LoadPatterns(listNotify, *importPeriod)
	notify.SetCarrier(carrier, *telegramGroup, *th)

	log.Println("connect to NATS server")
	opts := []nats.Option{nats.Name("NATS sniffer")}
	opts = mbroker.SetupConnOptions(opts)
	nc, err := nats.Connect(*natsURI, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	log.Println("subscribe to subject:", *subj)
	wg.Add(1)
	nc.Subscribe(*subj, func(m *nats.Msg) {
		var natsMessage mbroker.NatsMessage
		natsMessage.Subject = m.Subject
		natsMessage.Text = string(m.Data)
		isBanned := listRecords.CheckIfExists(natsMessage.Text)
		if !isBanned {
			repository.Call(natsMessage)
		} else {
			log.Println("skip banned message:", natsMessage.Text)
		}
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
