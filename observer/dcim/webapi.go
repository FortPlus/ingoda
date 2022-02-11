package dcim

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	BASE_URI = "/api/v1/"
)

type DeviceService struct {
	storage Storing
}

func NewDeviceService(storage Storing) DeviceService {
	return DeviceService{
		storage: storage,
	}
}

func (s *DeviceService) SetHandlers(r *mux.Router) {
	r.HandleFunc(BASE_URI, s.Get).Methods(http.MethodGet)
}

func (s *DeviceService) Get(w http.ResponseWriter, r *http.Request) {

	log.Println("DeviceService::Get get new request")

	// get query
	query, err := parseQuery(r)
	if err != nil {
		log.Println("DeviceService::Get failed parse json request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get result from repository
	devices, err := s.storage.Get(*query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// empty result
	if devices == nil {
		log.Println("DeviceService::Get no result for query:", query)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// serialize to json
	raw, err := createResponse(devices)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(raw)
}

func (s *DeviceService) Delete(id uint64) error {
	return nil
}

func (s *DeviceService) Add(device Device) error {
	return nil
}

func (s *DeviceService) CheckDNSName() error {
	return nil
}

func (s *DeviceService) CheckIsAlive() error {
	return nil
}
