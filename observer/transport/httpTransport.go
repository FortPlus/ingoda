package httpTransport

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fort.plus/fperror"
)

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

func PostJson(uri string, jsonPayload interface{}) error {
	var err error

	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(jsonPayload)
	log.Println("payload is:", payload)
	if err != nil {
		return fperror.Warning("Can't encode payload", err)
	}

	resp, err := http.Post(uri, "application/json", payload)
	defer resp.Body.Close()

	if err != nil {
		return fperror.Warning("Can't post data", err)
	}
    if resp.StatusCode != 200 {
        return fperror.Warning("got response with bad status" + resp.Status, nil)
    }
    log.Println("PostJson, response is", resp)
	return err
}
