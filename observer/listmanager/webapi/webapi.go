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

type listManager struct {
	Storage map[string]*banlist.ListRecords `json:"storage"`
	sync.RWMutex
}

func NewListManager() *listManager {
	lm := &listManager{
		Storage: make(map[string]*banlist.ListRecords),
	}
	go lm.clean()
	return lm
}

func (lm *listManager) Serialize() ([]byte, error) {
	data, err := json.Marshal(lm)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Deserialize(data []byte) (*listManager, error) {
	lm := NewListManager()
	if err := json.Unmarshal(data, &lm); err != nil {
		return nil, err
	}
	return lm, nil
}

func (lm *listManager) clean() {
	log.Println("listMapCleanEmpty()")
	for {
		time.Sleep(CLEANUP_PERIOD * time.Second)
		log.Println("listMapCleanEmpty:time to cleanup")
		lm.Lock()
		for key, val := range lm.Storage {
			if val.IsEmpty() {
				log.Println("listMapCleanEmpty:delete ", key)
				val.Close()
				delete(lm.Storage, key)
			}
		}
		lm.Unlock()
	}
}

func (lm *listManager) getOrCreateList(name string) *banlist.ListRecords {
	lm.RLock()
	defer lm.RUnlock()

	if _, exist := lm.Storage[name]; !exist {
		lm.Storage[name] = banlist.New(name)
	}
	return lm.Storage[name]
}

func (lm *listManager) getList(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.GetBannedList()")
	addHeaderParameters(w)
	params := mux.Vars(r)
	banList := lm.getOrCreateList(params["name"])

	records, err := banList.GetRecords()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(records)
}

func (lm *listManager) checkIfExist(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.CheckIfBanned")
	addHeaderParameters(w)
	params := mux.Vars(r)
	banList := lm.getOrCreateList(params["name"])

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

func (lm *listManager) addRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.AddRecord()")
	addHeaderParameters(w)
	params := mux.Vars(r)
	banList := lm.getOrCreateList(params["name"])

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

func (lm *listManager) deleteRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.DeleteRecord()")
	addHeaderParameters(w)
	params := mux.Vars(r)
	banList := lm.getOrCreateList(params["name"])

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

func SetHandlers(router *mux.Router) {

	lm := NewListManager()

	router.HandleFunc(BASE_URI, lm.getList).Methods(http.MethodGet)
	router.HandleFunc(CHECK_URI, lm.checkIfExist).Methods(http.MethodGet)
	router.HandleFunc(BASE_URI, lm.addRecord).Methods(http.MethodPost)
	router.HandleFunc(DELETE_URI, lm.deleteRecord).Methods(http.MethodDelete)
}

func addHeaderParameters(w http.ResponseWriter) {
	setCORS(w)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}
