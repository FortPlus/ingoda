package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"fort.plus/banlist"
	"fort.plus/config"
)

var (
	BAN_SERVER_PORT string = config.GetCurrent().BanServerPort
)

const (
	BASE_URL        = "/api"
	BANNED_LIST     = "/banned-list"
	CHECK_IF_BANNED = "/check-if-banned"
	API_V1          = "/v1"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc(BASE_URL+API_V1+BANNED_LIST, banlist.GetBannedList).Methods(http.MethodGet)
	router.HandleFunc(BASE_URL+API_V1+CHECK_IF_BANNED, banlist.CheckIfBanned).Methods(http.MethodGet)
	router.HandleFunc(BASE_URL+API_V1+BANNED_LIST, banlist.AddRecord).Methods(http.MethodPost)
	router.HandleFunc(BASE_URL+API_V1+BANNED_LIST+"/{id}", banlist.DeleteRecord).Methods(http.MethodDelete)

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(strings.Repeat("| ", len(ancestors)), t)
		return nil
	})

	http.ListenAndServe(BAN_SERVER_PORT, router)
}
