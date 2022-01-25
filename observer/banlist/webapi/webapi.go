package webapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"fort.plus/banlist"
)

const (
	API_URI         = "/api"
	BANNED_LIST     = "/banned-list"
	CHECK_IF_BANNED = "/exist"
	API_VERSION     = "/v1"

	BASE_URI = API_URI + API_VERSION + BANNED_LIST

	GET_LIST_URI = BASE_URI
	ADD_URI      = BASE_URI
	DELETE_URI   = BASE_URI + "/{id}"
	CHECK_URI    = BASE_URI + CHECK_IF_BANNED
)

var (
	banList *banlist.BannedRecords
)

func SetHandlers(router *mux.Router) {
	banList = banlist.NewBannedRecords("WebApi")

	router.HandleFunc(GET_LIST_URI, getBannedList).Methods(http.MethodGet)
	router.HandleFunc(CHECK_URI, checkIfBanned).Methods(http.MethodGet)
	router.HandleFunc(ADD_URI, addRecord).Methods(http.MethodPost)
	router.HandleFunc(DELETE_URI, deleteRecord).Methods(http.MethodDelete)
}

func getBannedList(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.GetBannedList()")
	addHeaderParameters(w)
	records, err := banList.GetRecords()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(records)
}

func checkIfBanned(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.CheckIfBanned")
	addHeaderParameters(w)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	isExists := banList.CheckIfExists(string(b[:]))
	if isExists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func addRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.AddRecord()")
	addHeaderParameters(w)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	item := new(banlist.Item)
	err = json.Unmarshal(b, &item)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	banList.AddRecord(*item)
	w.WriteHeader(http.StatusCreated)
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.DeleteRecord()")
	addHeaderParameters(w)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = banList.Delete(uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addHeaderParameters(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

}
