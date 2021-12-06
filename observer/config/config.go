//
// singleton with configuration data
//
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
	"sync"
)

const(
    CONFIG_FILE_ENV_VAR = "FP_LOG_CONF"
)

type config struct {
    SyslogSocket string  `json:"syslog.socket"`

	TelegramToken string `json:"telegram.token"`
	TelegramNotificationGroup int64 `json:"telegram.notification.group"`
	TelegramNotificationDiffThreshold int `json:"telegram.notification.diff.threshold"`


	NatsConnectionString string `json:"nats.connection.string"`
	NatsSyslogSubject string `json:"nats.syslog.subject"`


	ExtSnifferNotifyTelegram []string `json:"extension.sniffer.notify.telegram"`

	empty bool
}

var once sync.Once
var currentConfig *config


// Method to get configuration structure
// return immutable structure by value
func GetCurrent() config {
    once.Do(func() {
        confFile := lookupConfigFilePathEnv()
        currentConfig = readFromFile(confFile)
    })
    return *currentConfig
}

func lookupConfigFilePathEnv() string {
    configFile, ok := os.LookupEnv(CONFIG_FILE_ENV_VAR)
    if !ok {
        log.Panicf("please set environment variable %s with path to config file", CONFIG_FILE_ENV_VAR)
    }
    return configFile
}

func readFromFile(path string) *config{
    currentConfig = &config{empty: true}
    file, err := ioutil.ReadFile(path)
    if err != nil {
        log.Panicf("can't read configuration file:%s, %s", path, err)
    }
    err = json.Unmarshal([]byte(file), &currentConfig)
    if err != nil {
            log.Panicf("configuration file:%s, has an error: %s", path, err)
    }
    currentConfig.empty = false
    return currentConfig
}
