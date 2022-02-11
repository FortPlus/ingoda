package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// init components
	var repository dcim.Storing = dcim.NewRepoInMemory(*path)
	var service dcim.DeviceService = dcim.NewDeviceService(repository)

	if err := repository.Initialize(); err != nil {
		log.Fatal("main:: ", err)
	}

	// setup handlers
	router := mux.NewRouter()
	service.SetHandlers(router)

	// run http server
	if err := http.ListenAndServe(*portNum, router); err != nil {
		log.Fatal(err)
	}
}
