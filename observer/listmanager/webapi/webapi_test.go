package webapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	banlist "fort.plus/listmanager"
	"github.com/gorilla/mux"
)

func TestAddRecord(t *testing.T) {

	var data []byte = []byte("{\"pattern\":\"test\", \"expired_at\":\"2022-02-24T16:15:02.296629921+03:00\"}")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/list1/", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	lm := NewListManager()
	lm.SetHandlers(mux.NewRouter())
	lm.addRecord(w, req)
	if w.Code != http.StatusCreated {
		t.Fatal("bad status code", w.Code)
	}
	fmt.Println(lm.Storage)
	for k, v := range lm.Storage {
		fmt.Println(k, v)
	}
}

func TestGetList(t *testing.T) {
	lm := NewListManager()
	lm.SetHandlers(mux.NewRouter())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/list1/", nil)
	w := httptest.NewRecorder()
	lm.getList(w, req)

	data, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(data))
}

func TestSerialize(t *testing.T) {

	lm := NewListManager()
	fmt.Println(lm)

	lm.Storage["list1"] = banlist.New("list1")
	list1 := lm.Storage["list1"]
	list1.AddRecord(banlist.Item{Pattern: "pattern1", ExpiredAt: time.Now().Add(time.Minute * 2)})
	list1.AddRecord(banlist.Item{Pattern: "pattern2", ExpiredAt: time.Now().Add(time.Minute * 3)})

	list2 := banlist.New("list2")
	lm.Storage["list2"] = list2
	list2.AddRecord(banlist.Item{Pattern: "pattern1", ExpiredAt: time.Now().Add(time.Minute * 4)})
	list2.AddRecord(banlist.Item{Pattern: "pattern2", ExpiredAt: time.Now().Add(time.Minute * 5)})

	t.Run("Serialze", func(t *testing.T) {

		data, err := lm.Serialize()
		fmt.Println(err)
		fmt.Println(string(data))
	})

	t.Run("Deserialize", func(t *testing.T) {
		data, err := lm.Serialize()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("\n\n---Unmarshal")
		lm2, err := Deserialize(data)
		if err != nil {
			t.Fatal(err)
		}
		// json.Unmarshal(data, &lm2)
		fmt.Println(lm2)

		for _, value := range lm2.Storage {
			// value.GetPatterns()
			fmt.Println("pattern", value.GetPatterns())
		}
	})
}

func TestTime(t *testing.T) {
	t1 := time.Now()
	tBefore := t1.Add(time.Hour * -1)
	tAfter := t1.Add(time.Hour + 1)
	fmt.Println(t1.After(tBefore))
	fmt.Println(t1.Before(tAfter))
}
