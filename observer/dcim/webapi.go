package dcim

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"fort.plus/fperror"
	"github.com/gorilla/mux"
)

const (
	BASE_URI  = "/api/v1/"
	WHOIS_URI = "/api/v1/whois"
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
	devices, err := s.storage.Get(query)
	if err != nil {
		log.Println(err)
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

// WhoisService represent whois information endpoint
type WhoisService struct {
	data []byte
}

func NewWhoisService() WhoisService {
	return WhoisService{}
}

func (wh *WhoisService) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		err = fperror.Warning("can't open file", err)
		log.Println(err)
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		err = fperror.Warning("can't read all bytes", err)
		return err
	}
	wh.data = bytes
	return nil
}

func (wh *WhoisService) match(pattern string) []string {
	var result []string
	re := regexp.MustCompile(pattern)
	scanner := bufio.NewScanner(bytes.NewReader(wh.data))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if re.Match([]byte(line)) {
			result = append(result, line)
		}
	}
	return result
}

func (wh *WhoisService) SetHandlers(r *mux.Router) {
	r.HandleFunc(WHOIS_URI, wh.Get).Methods(http.MethodGet)
}

func (wh *WhoisService) Get(w http.ResponseWriter, r *http.Request) {
	// parse query
	var q struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		log.Println("WhoisService::Get failed parse json request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// grep result
	result := wh.match(q.Query)

	// write result
	raw, err := json.Marshal(result)
	if err != nil {
		log.Println("WhoisService::get failed marshal response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(raw)
}
