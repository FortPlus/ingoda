// Provide helper functions to work with NATS message queue server
package mbroker

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"hash/crc32"
	"math/rand"

	"github.com/nats-io/nats.go"

	"fort.plus/config"
	"fort.plus/fperror"
)

const (
	MAX_BATCH_SALT_VALUE = 1000
)

var wg sync.WaitGroup
var nc *nats.Conn

var NATS_SERVER string = config.GetCurrent().NatsConnectionString

type CallbackFunction = func(response string, err error)

func natsConnect() {
	var err error = nil
	if nc != nil && nc.IsConnected() {
		return
	}

	opts := []nats.Option{nats.Name("wait-for-subject")}
	opts = setupConnOptions(opts)

	// Connect to a server
	nc, err = nats.Connect(NATS_SERVER, opts...)
	if err != nil {
		log.Fatal(fperror.Warning("can't connect to NATS server", err))
	}
}

//
//  Wait for message with specific subject
//
func WaitForSubject(subject string, timeout int, callback CallbackFunction) (string, error) {
	response := "{}"
	var err error = nil

	natsConnect()

	wg.Add(1)
	nc.Subscribe(subject, func(m *nats.Msg) {
		response = string(m.Data)
		wg.Done()
	})
	nc.Flush()

	if err = nc.LastError(); err != nil {
		log.Println(err)

	}

	if waitTimeout(&wg, time.Duration(timeout)*time.Second) {
		log.Println("Timed out waiting for wait group")
		err = errors.New("Timeout error")
	} else {
		log.Println("Wait group finished")
	}

	callback(response, err)

	return response, err
}

//
// Send request message to NATS
//
func SendRequest(subject string, data string, timeout int) (string, error) {
	response := ""
	natsConnect()

	msg, err := nc.Request(subject, []byte(data), time.Duration(timeout)*time.Second)
	if err != nil {
		log.Println(err)
	} else {
		response = string(msg.Data)
	}
	return response, err
}

//
// Close connection to NATS server
//
func Close() {
	if nc != nil {
		nc.Close()
	}
}

//
// Publish message to NATS
//
func Publish(subject string, data string) error {
	natsConnect()
	err := nc.Publish(subject, []byte(data))
	if err != nil {
		log.Println(fperror.Warning("can't publish data,", err))
	}
	return err
}

//
//  Calculate batch number using combination of crc32(TXT + RANDOM INT)
//
func GetBatchId(data string) string {
	batchString := fmt.Sprintf("%s.%d", data, rand.Intn(MAX_BATCH_SALT_VALUE))
	batchId := crc32.ChecksumIEEE([]byte(batchString))
	return strconv.FormatUint(uint64(batchId), 10)
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
		log.Printf("Exiting, no servers available")
	}))
	return opts
}
