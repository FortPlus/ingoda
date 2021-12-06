package main

import (
    "log"
    "sync"
    "time"

    "github.com/nats-io/nats.go"

    _ "fort.plus/extension/sniffer"
    "fort.plus/repository"
    "fort.plus/config"
    "fort.plus/mbroker"
)

var (
    NATS_SUBJECT string = config.GetCurrent().NatsSyslogSubject
    NATS_SERVER string = config.GetCurrent().NatsConnectionString
)

var wg sync.WaitGroup

func main() {

    opts := []nats.Option{nats.Name("NATS sniffer")}
    opts = setupConnOptions(opts)
    nc, err := nats.Connect(NATS_SERVER, opts...)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    wg.Add(1)

    nc.Subscribe(NATS_SUBJECT+".>", func(m *nats.Msg) {
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

