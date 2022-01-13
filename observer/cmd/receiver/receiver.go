package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"

	"fort.plus/config"
	"fort.plus/fperror"
	"fort.plus/mbroker"
)

const (
	MAX_UDP_SIZE = math.MaxUint16
)

var (
	RECEIVER_UDP_SOCKET string = config.GetCurrent().ReceiverUdpSocket
	NATS_SUBJECT        string = config.GetCurrent().NatsSyslogSubject
)

func main() {
	buff := make([]byte, MAX_UDP_SIZE)
	var message mbroker.SyslogMessage

	conn := prepareSocket(RECEIVER_UDP_SOCKET)
	defer conn.Close()
	defer mbroker.Close()

	for {
		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			log.Fatal(fperror.Critical("Can't read from unix socket", err))
		}
		log.Printf("n: %d, addr: %v, err: %v, buf: %s\n", n, addr, err, buff[0:n])

		err = json.Unmarshal(buff[0:n], &message)
		if err != nil {
			log.Println(fperror.Warning("can't unmarshall message", err))
			continue
		}
		subj := fmt.Sprintf("%s.%s.%s.%s", NATS_SUBJECT, message.Facility, message.Severity, message.Host)

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
