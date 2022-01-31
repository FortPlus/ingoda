package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"gopkg.in/alecthomas/kingpin.v2"

	"fort.plus/im"
	"fort.plus/im/telegram"
)

const (
	MAX_DELIVERY_TIME         time.Duration = 10 //seconds
	SLEEP_BEFORE_SEND_MESSAGE time.Duration = 3  //seconds
)

var (
	wg            sync.WaitGroup
	app                      = kingpin.New(filepath.Base(os.Args[0]), "Check syslog workflow")
	token                    = app.Flag("telegram.token", "Telegram token.").Required().OverrideDefaultFromEnvar("TELEGRAM_TOKEN").String()
	telegramGroup            = app.Flag("telegram.group", "Telegram notification group.").Required().OverrideDefaultFromEnvar("TELEGRAM_GROUP").String()
	natsURI                  = app.Flag("nats.uri", "URI of the NATS server.").Default("nats://localhost:4222").OverrideDefaultFromEnvar("NATS_URI").String()
	syslogSubj               = app.Flag("syslog.subj", "Syslog messsages subject prefix in nats.").Default("slg").OverrideDefaultFromEnvar("SYSLOG_SUBJ").String()
	carrier       im.Carrier = telegram.New(*token)
)

func main() {
	log.SetOutput(os.Stdout)

	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	timestampInt := time.Now().UTC().UnixNano()
	timestamp := strconv.FormatInt(timestampInt, 10)
	messageToSend := fmt.Sprintf("%s, One, Two, Freddy's Coming For You", timestamp)

	go listenNats(messageToSend)
	time.Sleep(SLEEP_BEFORE_SEND_MESSAGE * time.Second)
	syslogSend(messageToSend)
	time.Sleep(MAX_DELIVERY_TIME * time.Second)

}

func syslogSend(message string) {
	sysLog, err := syslog.Dial("udp", "localhost:514", syslog.LOG_WARNING|syslog.LOG_DAEMON, "selfcheck")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(sysLog, message)
}

// TODO: replace with mbroker methods
func listenNats(message string) bool {
	result := false

	opts := []nats.Option{nats.Name("selfcheck-syslog")}
	opts = setupConnOptions(opts)

	// Connect to a server
	nc, err := nats.Connect(*natsURI, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	wg.Add(1)

	nc.Subscribe(*syslogSubj+".>", func(m *nats.Msg) {
		if strings.Contains(string(m.Data), message) {
			fmt.Printf("got the same string [%s]: '%s'\n", m.Subject, string(m.Data))
			wg.Done()
		} else {
			fmt.Printf("skip received on [%s]: '%s'\n", m.Subject, string(m.Data))
		}
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	if waitTimeout(&wg, (MAX_DELIVERY_TIME+SLEEP_BEFORE_SEND_MESSAGE)*time.Second) {
		carrier.Send(*telegramGroup, "Dear engineers, syslog doesn't work correctly")
		log.Println("Timed out waiting for wait group")
	} else {
		log.Println("Wait group finished")
	}

	return result
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
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
