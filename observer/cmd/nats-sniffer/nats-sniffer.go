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
	_ "fort.plus/extension/sniffer/notify"
	"fort.plus/im/telegram"
	"fort.plus/repository"

	"fort.plus/mbroker"
)

var (
	app           = kingpin.New(filepath.Base(os.Args[0]), "Nats message sniffer")
	token         = app.Flag("telegram.token", "Telegram token.").Required().OverrideDefaultFromEnvar("TELEGRAM_TOKEN").String()
	telegramGroup = app.Flag("telegram.group", "Telegram notification group.").Required().OverrideDefaultFromEnvar("TELEGRAM_GROUP").String()
	th            = app.Flag("ld.th", "Levenshtein distance threshold.").Default("14").OverrideDefaultFromEnvar("LD_TH").Int()
	natsURI       = app.Flag("nats.uri", "URI of the NATS server.").Default("nats://localhost:4222").OverrideDefaultFromEnvar("NATS_URI").String()
	subj          = app.Flag("nats.subj", "NATS message subject").Default("slg.>").OverrideDefaultFromEnvar("NATS_SUBJ").String()
)

var wg sync.WaitGroup

func main() {
	log.SetOutput(os.Stdout)

	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	carrier := telegram.New(*token)
	notify.Register(carrier, *telegramGroup, *th)

	opts := []nats.Option{nats.Name("NATS sniffer")}
	opts = setupConnOptions(opts)
	nc, err := nats.Connect(*natsURI, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	wg.Add(1)

	nc.Subscribe(*subj+".>", func(m *nats.Msg) {
		var natsMessage mbroker.NatsMessage
		natsMessage.Subject = m.Subject
		natsMessage.Text = string(m.Data)
		repository.Call(natsMessage)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 1 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		log.Printf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatal("Exiting, no servers available")
	}))
	return opts
}
