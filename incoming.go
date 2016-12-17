package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stvp/rollbar"
)

func handleIncoming(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Printf("incoming request %s", string(reqB))

	incomingMessage := IncomingMessage{}
	err = json.Unmarshal(reqB, &incomingMessage)
	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if incomingMessage.Object != "page" {
		http.Error(w, "only allowed to chat via fb page", http.StatusForbidden)
		return
	}

	if incomingMessage.Entries == nil {
		rollbar.Error(rollbar.ERR, fmt.Errorf("entry is not provided"))
		http.Error(w, "", http.StatusExpectationFailed)
		return
	}

	for _, entry := range *incomingMessage.Entries {
		for _, message := range entry.Messaging {
			if message.MessageData.Text == "" {
				continue
			}
			log.Println("handling outgoing message - ", message.MessageData.Text)
			handleOutgoing(w, message)
		}
	}

	return
}
