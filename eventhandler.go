package eventhandler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// HeaderTokens contains general information about the event and supervisor
// Adopted from http://supervisord.org/events.html#header-tokens
type HeaderTokens struct {
	Ver        string
	Server     string
	Serial     string
	Pool       string
	PoolSerial string
	EventName  string
	len        int
}

// ValidEventNames contains the events supported by supervisor.
// Adopted from http://supervisord.org/events.html#event-types
var ValidEventNames = []string{
	"EVENT",
	"PROCESS_STATE",
	"PROCESS_STATE_STARTING",
	"PROCESS_STATE_RUNNING",
	"PROCESS_STATE_BACKOFF",
	"PROCESS_STATE_EXITED",
	"PROCESS_STATE_STOPPED",
	"PROCESS_STATE_FATAL",
	"PROCESS_STATE_UNKNOWN",
	"REMOTE_COMMUNICATION",
	"PROCESS_LOG",
	"PROCESS_LOG_STDOUT",
	"PROCESS_LOG_STDERR",
	"PROCESS_COMMUNICATION",
	"PROCESS_COMMUNICATION_STDOUT",
	"PROCESS_COMMUNICATION_STDERR",
	"SUPERVISOR_STATE_CHANGE",
	"SUPERVISOR_STATE_CHANGE_RUNNING",
	"SUPERVISOR_STATE_CHANGE_STOPPING",
	"TICK",
	"TICK_5",
	"TICK_60",
	"TICK_3600",
	"PROCESS_GROUP",
	"PROCESS_GROUP_ADDED",
	"PROCESS_GROUP_REMOVED",
}

// EventHandler is the main service struct to this package.
// It should be initialized with HandleEvent function and then started.
type EventHandler struct {
	processors map[string]EventProcessor
}

// EventProcessor defines the actual event processing function
// This should be provided by the client
type EventProcessor func(HeaderTokens, map[string]string)

// New creates, initalizes, and returns EventHandler
func New() *EventHandler {
	return &EventHandler{make(map[string]EventProcessor)}
}

// HandleEvent puts a new processor to the EventHandler, which will
// be used while processing supervisord events.
func (h *EventHandler) HandleEvent(eventName string, processor EventProcessor) error {
	valid := false
	for _, n := range ValidEventNames {
		if n == eventName {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("%s is not a valid event! Valid events are: %v", eventName, ValidEventNames)
	}

	if h.processors == nil {
		h.processors = make(map[string]EventProcessor)
	}
	h.processors[eventName] = processor
	return nil
}

// Start is the blocking event handling function for EventHandler
// Should be called as last step
func (h *EventHandler) Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("READY")

		header, dataMap, err := h.readHeaderAndPayload(reader)

		if err != nil {
			log.Printf("Processing event failed, probably not your fault, error: %s\n", err)
		}
		go h.processEvent(header, dataMap)

		fmt.Print("RESULT 2\nOK")
	}
}

// Reads header tokens and payload from reader
// Returns parsed header tokens and payload
func (h *EventHandler) readHeaderAndPayload(reader *bufio.Reader) (headerTokens HeaderTokens, payloadMap map[string]string, err error) {
	headerLine, err := reader.ReadString('\n')
	if err != nil {
		err = errors.Wrap(err, "Reading header line failed")
		return
	}

	headerTokens = h.parseHeaderTokens(headerLine)

	payload := make([]byte, headerTokens.len)
	_, err = reader.Read(payload)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Reading payload line failed, Headers: %v\n", headerTokens))
		return
	}

	payloadMap = h.parseTokensToMap(string(payload))

	return
}

func (h *EventHandler) processEvent(header HeaderTokens, payload map[string]string) {
	processor, ok := h.processors[header.EventName]
	if !ok {
		// for generic event types like PROCESS_STATE
		for event, p := range h.processors {
			if strings.HasPrefix(header.EventName, event) {
				processor = p
			}
		}
	}
	if processor == nil {
		return
	}

	processor(header, payload)
}

// Receives space separated {key}:{value} string pairs,
// creates a map where key -> value
func (h *EventHandler) parseTokensToMap(tokens string) map[string]string {
	tokenMap := make(map[string]string)
	tokenList := strings.Split(strings.TrimSpace(tokens), " ")
	for _, entry := range tokenList {
		splited := strings.Split(entry, ":")
		tokenMap[splited[0]] = splited[1]
	}
	return tokenMap
}

// Parses given header string, extracts values & returns HeaderTokens
// Example header tokens:
// ver:3.0 server:supervisor serial:21 pool:listener poolserial:10 eventname:PROCESS_COMMUNICATION_STDOUT len:54
func (h *EventHandler) parseHeaderTokens(header string) HeaderTokens {
	headerMap := h.parseTokensToMap(header)
	len, err := strconv.Atoi(headerMap["len"])
	if err != nil {
		log.Fatal(err)
	}

	return HeaderTokens{
		headerMap["ver"],
		headerMap["server"],
		headerMap["serial"],
		headerMap["pool"],
		headerMap["poolserial"],
		headerMap["eventname"],
		len,
	}
}
