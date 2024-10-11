package services

import (
	"fmt"
	"log"
	"net/http"
)

func SendRemoteCommand(listingAddr string, command string) bool {
	resp, err := http.Get("http://" + listingAddr + "/commands/" + command)
	if err != nil {
		log.Panic(err)
	}
	log.Print(resp.StatusCode)
	return (resp.StatusCode == 200)
}

func IsRemoteCommandAvaialable(listingAddr string) bool {
	resp, err := http.Get("http://" + listingAddr + "/status")
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}

func StartRemoteCommandListing(listingAddr string,
	errorCallback func(err error),
	onCommandCallback func(command string) bool) {

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ready")
	})

	http.HandleFunc("/commands/showHide", func(w http.ResponseWriter, r *http.Request) {
		res := onCommandCallback("showHide")
		if res {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	go func() {
		err := http.ListenAndServe(listingAddr, nil)
		if err != nil {
			errorCallback(err)
		}
	}()
}
