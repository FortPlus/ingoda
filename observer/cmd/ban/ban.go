package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"fort.plus/banlist/webapi"
	"fort.plus/config"
)

var (
	BAN_SERVER_PORT string = config.GetCurrent().BanServerPort
)

func main() {
	router := mux.NewRouter()
	webapi.SetHandlers(router)
	http.ListenAndServe(BAN_SERVER_PORT, router)
}
