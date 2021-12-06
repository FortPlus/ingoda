package main

import (
    "net"
    "os"
    "log"
    "time"
    "strings"
    "encoding/json"

    "github.com/nats-io/nats.go"

    "fort.plus/config"
    "fort.plus/mbroker"
    "fort.plus/fperror"
)

var (
    SYSLOG_SOCKET string = config.GetCurrent().SyslogSocket
    NATS_SERVER string = config.GetCurrent().NatsConnectionString
    NATS_SUBJECT string = config.GetCurrent().NatsSyslogSubject
)

func main() {

    conn := prepareUnixSocketConnection(SYSLOG_SOCKET)
 	defer conn.Close()

    nc := prepareNatsConnection(NATS_SERVER)
    defer nc.Close()
   
    buff := make([]byte, 1024)
    var msg mbroker.SyslogMessage
    var sb strings.Builder
    for {
        sb.Reset()
        n, addr, err := conn.ReadFrom(buff)
        if err!= nil {
            log.Fatal(fperror.Critical("Can't read from unix socket", err))
        }

        log.Printf("n: %d, addr: %v, err: %v, buf: %s\n", n, addr, err, buff[0:n])
        
        json.Unmarshal(buff[0:n], &msg)

        sb.WriteString(NATS_SUBJECT)
        sb.WriteString(".")
        sb.WriteString(msg.Facility)
        sb.WriteString(".")
        sb.WriteString(msg.Severity)
        sb.WriteString(".")
        sb.WriteString(msg.Host)

        err = nc.Publish(sb.String(), buff[0:n])
        if err != nil{
          log.Fatal(fperror.Critical("Can't publish", err))
        }
    }
}


func prepareNatsConnection(natsServer string) (*nats.Conn) {
    opts := []nats.Option{nats.Name("NATS sniffer")}
    opts = setupConnOptions(opts)
    nc, err := nats.Connect(natsServer, opts...)
    if err != nil {
        log.Fatal(fperror.Critical("Can't connect to NATS server", err))
    }
    return nc
}

func prepareUnixSocketConnection(socketName string) (*net.UnixConn) {
    os.Remove(socketName)
    socket, err := net.ResolveUnixAddr("unixgram", socketName)
    if err != nil {
        log.Fatal(fperror.Critical("Can't resolve unix addr", err))
    }
    conn, err := net.ListenUnixgram("unixgram", socket)
    if err != nil {
        log.Fatal(fperror.Critical("Can't listen unix domain socket", err))
    }
    return conn
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

