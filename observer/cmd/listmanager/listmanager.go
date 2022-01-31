package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"

	"fort.plus/listmanager/webapi"
)

func main() {
	log.SetOutput(os.Stdout)

	app := kingpin.New(filepath.Base(os.Args[0]), "List manager")
	portNum := app.Flag("port", "Port number to listen.").Default(":9190").OverrideDefaultFromEnvar("LIST_MANAGER_PORT").String()
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	router := mux.NewRouter()
	webapi.SetHandlers(router)
	http.ListenAndServe(*portNum, router)
}
