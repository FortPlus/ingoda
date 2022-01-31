package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"

	"fort.plus/fperror"
	"fort.plus/mbroker"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	MAX_UDP_SIZE = math.MaxUint16
)

var (
	app     = kingpin.New(filepath.Base(os.Args[0]), "syslog to NATS message converter")
	portUDP = app.Flag("port", "UDP port number to listen.").Default(":9191").OverrideDefaultFromEnvar("LIST_MANAGER_PORT").String()
	natsURI = app.Flag("nats.uri", "NATS server URI").Default("nats://localhost:4222").OverrideDefaultFromEnvar("NATS_URI").String()
	subj    = app.Flag("subj.prefix", "Subject prefix for syslog messages.").Default("slg").OverrideDefaultFromEnvar("SUBJ_PREFIX").String()
)

func main() {
	log.SetOutput(os.Stdout)

	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	buff := make([]byte, MAX_UDP_SIZE)
	var message mbroker.SyslogMessage

	log.Println("listen socket", *portUDP)
	conn := prepareSocket(*portUDP)

	log.Println("connect to NATS:", *natsURI)
	mbroker.New(*natsURI)
	defer conn.Close()
	defer mbroker.Close()
	log.Println("ready to resend messages")

	for {
		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			log.Fatal(fperror.Critical("Can't read from socket", err))
		}
		log.Printf("n: %d, addr: %v, err: %v, buf: %s\n", n, addr, err, buff[0:n])

		err = json.Unmarshal(buff[0:n], &message)
		if err != nil {
			log.Println(fperror.Warning("can't unmarshall message", err))
			continue
		}
		subj := fmt.Sprintf("%s.%s.%s.%s", *subj, message.Facility, message.Severity, message.Host)

		err = mbroker.Publish(subj, string(buff[0:n]))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func prepareSocket(socket string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp4", socket)
	if err != nil {
		err = fperror.Critical("can't resolve udp address", err)
		log.Fatal(err)
	}
	con, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		err = fperror.Critical("can't listen udp socket", err)
		log.Fatal(err)
	}
	return con
}
