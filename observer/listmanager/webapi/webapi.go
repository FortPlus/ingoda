package webapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"

	banlist "fort.plus/listmanager"
)

const (
	BASE_URI       = "/api/v1/{name}"
	CHECK_URI      = BASE_URI + "/exist"
	DELETE_URI     = BASE_URI + "/{id}"
	CLEANUP_PERIOD = 180
)

var (
	listMap map[string]*banlist.ListRecords = map[string]*banlist.ListRecords{}
	lock                                    = sync.RWMutex{}
)

func listMapCleanEmpty() {
	log.Println("listMapCleanEmpty()")
	for {
		time.Sleep(CLEANUP_PERIOD * time.Second)
		log.Println("listMapCleanEmpty:time to cleanup")
		lock.Lock()
		for key, val := range listMap {
			if val.IsEmpty() {
				log.Println("listMapCleanEmpty:delete ", key)
				val.Close()
				delete(listMap, key)
			}
		}
		lock.Unlock()
	}
}

func getOrCreateList(name string) *banlist.ListRecords {
	lock.RLock()
	defer lock.RUnlock()

	if _, exist := listMap[name]; !exist {
		listMap[name] = banlist.New(name)
	}
	return listMap[name]
}

func SetHandlers(router *mux.Router) {
	go listMapCleanEmpty()

	router.HandleFunc(BASE_URI, getList).Methods(http.MethodGet)
	router.HandleFunc(CHECK_URI, checkIfExist).Methods(http.MethodGet)
	router.HandleFunc(BASE_URI, addRecord).Methods(http.MethodPost)
	router.HandleFunc(DELETE_URI, deleteRecord).Methods(http.MethodDelete)
}

func getList(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.GetBannedList()")
	addHeaderParameters(w)
	params := mux.Vars(r)
	banList := getOrCreateList(params["name"])

	records, err := banList.GetRecords()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(records)
}

func checkIfExist(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.CheckIfBanned")
	addHeaderParameters(w)
	params := mux.Vars(r)
	banList := getOrCreateList(params["name"])

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
	params := mux.Vars(r)
	banList := getOrCreateList(params["name"])

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
	banList := getOrCreateList(params["name"])

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
	setCORS(w)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}
