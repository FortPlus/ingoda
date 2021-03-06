package httpTransport

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"fort.plus/fperror"
)

type JsonConverter interface {
	UnmarshallJSON([]byte) error
	MarshallJSON(data []byte) ([]byte, error)
}

func GetAndUnmarshall(uri string, response interface{}) error {
	var err error
	resp, err := http.Get(uri)
	if err != nil {
		return fperror.Warning("Can't get data", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fperror.Warning("Can't read from channel", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return fperror.Warning("Can`t unmarshall json", err)
	}
	return err
}

func Get(uri string, f JsonConverter) error {
	var err error

	resp, err := http.Get(uri)
	if err != nil {
		return fperror.Warning("Can't get data", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fperror.Warning("Can't read from channel", err)
	}
	err = f.UnmarshallJSON(body)
	if err != nil {
		return fperror.Warning("Can`t unmarshall json", err)
	}
	return err
}

func PostJson(uri string, jsonPayload interface{}) error {
	var err error

	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(jsonPayload)
	log.Println("payload is:", payload)
	if err != nil {
		return fperror.Warning("Can't encode payload", err)
	}

	resp, err := http.Post(uri, "application/json", payload)
	if err != nil {
		return fperror.Warning("Can't post data", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fperror.Warning("got response with bad status"+resp.Status, nil)
	}
	log.Println("PostJson, response is", resp)
	return err
}

func Request(method string, uri string, body io.Reader) (*http.Response, error) {
	var err error
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	return resp, err
}

func Delete(uri string) error {
	var err error
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(respBody)
	return err
}
