package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"

	"fort.plus/dcim"
)

func main() {

	log.SetOutput(os.Stdout)

	// parse flags
	app := kingpin.New(filepath.Base(os.Args[0]), "DCIM")
	portNum := app.Flag("port", "Port number to listen.").Default(":38000").OverrideDefaultFromEnvar("DCIM_PORT").String()
	path := app.Flag("db", "DCIM DB file.").Default("/var/spool/dcim-db.json").OverrideDefaultFromEnvar("DCIM_DB").String()
	whoisFileName := app.Flag("whois.file", "Path to file with whois database.").Default("/var/spool/whois.txt").OverrideDefaultFromEnvar("WHOIS_FILE").String()
	whoisUpdateTimer := app.Flag("whois.timer", "Update timer.").Default("30m").OverrideDefaultFromEnvar("WHOIS_UPDATE_TIMER").Duration()
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// init components
	var repository dcim.Storing = dcim.NewRepoInMemory(*path)
	var service *dcim.DeviceService = dcim.NewDeviceService(repository)
	var whois *dcim.WhoisService = dcim.NewWhoisService(*whoisFileName, *whoisUpdateTimer)

	if err := repository.Initialize(); err != nil {
		log.Fatal("main:: ", err)
	}

	// setup handlers
	router := mux.NewRouter()
	dcim.SetHandlers(router, service, whois)

	// run http server
	go func() {
		if err := http.ListenAndServe(*portNum, router); err != nil {
			log.Fatal(err)
		}
	}()

	// gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ctx.Done():
		log.Println("Services stops by IPC signal")
		// TODO: implement service graceful shutdown
		// device-service, whois-service

		stop()
	}
}
