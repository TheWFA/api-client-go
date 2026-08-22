// Command webhook demonstrates verifying and handling an inbound WFA webhook
// delivery over HTTP.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TheWFA/api-client-go/webhooks"
)

func main() {
	publicKey := os.Getenv("WFA_WEBHOOK_PUBLIC_KEY")

	http.HandleFunc("/webhooks/wfa", func(w http.ResponseWriter, r *http.Request) {
		event, err := webhooks.ConstructEvent(r, publicKey, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch e := event.(type) {
		case webhooks.GoalScoredEvent:
			fmt.Printf("goal scored at match time %v\n", e.MatchTime)
		case webhooks.CardIssuedEvent:
			fmt.Printf("card issued: %s\n", e.CardType)
		case webhooks.MatchStatusChangedEvent:
			fmt.Printf("status changed: %s -> %s\n", e.PreviousStatus, e.NewStatus)
		case webhooks.PingEvent:
			fmt.Println("received a ping")
		default:
			fmt.Printf("received %s\n", event.EventType())
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
