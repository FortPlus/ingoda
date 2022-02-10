package main

import (
	"fmt"
	"log"
	"net/http"

	"fort.plus/dcim"
)

func main() {

	// init components
	var repository dcim.DeviceRepository = dcim.NewRepoInMemory()
	var service dcim.DeviceService = dcim.NewDeviceService(repository)
	fmt.Println(service)

	if err := repository.Initialize(); err != nil {
		log.Fatal("main:: ", err)
	}

	// setup handlers

	// run http server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
