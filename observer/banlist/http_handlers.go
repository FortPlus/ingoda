package banlist

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	banList *BannedRecords = NewBannedRecords("ban list")
)

func GetBannedList(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.GetBannedList()")
	records, err := banList.GetRecords()
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(records)
}

func CheckIfBanned(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.CheckIfBanned()")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	isExists := banList.CheckIfExists(string(b[:]))
	w.Header().Add("Content-Type", "application/json")
	if isExists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func AddRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("banlist.AddRecord()")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	item := new(Item)
	err = json.Unmarshal(b, &item)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	banList.AddRecord(*item)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
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
