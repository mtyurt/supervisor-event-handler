package main

import (
	"fmt"
	"os"

	eventhandler "github.com/mtyurt/supervisor-event-handler"
)

func main() {
	handler := eventhandler.EventHandler{}

	handler.HandleEvent("PROCESS_STATE", func(header eventhandler.HeaderTokens, payload map[string]string) {

		fmt.Fprintf(os.Stderr, "event: %s, payload: %v\n", header.EventName, payload)

	})
	handler.Start()
}
