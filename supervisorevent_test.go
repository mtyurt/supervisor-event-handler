package supervisorevent

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseHeader(t *testing.T) {
	h := EventHandler{}
	var tests = []struct {
		expected HeaderTokens
		given    string
	}{
		{
			HeaderTokens{
				Ver:        "3.0",
				Server:     "supervisor",
				Serial:     "21",
				Pool:       "listener",
				PoolSerial: "10",
				EventName:  "PROCESS_COMMUNICATION_STDOUT",
				len:        54,
			},
			"ver:3.0 server:supervisor serial:21 pool:listener poolserial:10 eventname:PROCESS_COMMUNICATION_STDOUT len:54",
		},
		{
			HeaderTokens{
				Ver:        "3.0",
				Server:     "supervisor",
				Serial:     "21",
				Pool:       "listener",
				PoolSerial: "12",
				EventName:  "PROCESS_STATE_STOPPED",
				len:        50,
			},
			"len:50 serial:21 ver:3.0 pool:listener poolserial:12 server:supervisor eventname:PROCESS_STATE_STOPPED",
		},
	}
	for _, tt := range tests {
		actual := h.parseHeaderTokens(tt.given)
		if actual != tt.expected {
			t.Errorf("parseHeaderTokens(%s): expected %v, actual %v", tt.given, tt.expected, actual)
		}
	}
}

func TestRegisterEventProcessor(t *testing.T) {
	h := EventHandler{}
	var tests = []struct {
		expected error
		given    struct {
			event   string
			handler EventProcessor
		}
	}{
		{nil, struct {
			event   string
			handler EventProcessor
		}{
			"PROCESS_STATE",
			nil,
		}},
		{errors.New(fmt.Sprintf("invalidevent is not a valid event! Valid events are: %v", VALID_EVENT_NAMES)), struct {
			event   string
			handler EventProcessor
		}{
			"invalidevent",
			nil,
		}},
	}
	for _, tt := range tests {
		actual := h.RegisterEventProcessor(tt.given.event, tt.given.handler)
		if (actual != nil && tt.expected != nil) && actual.Error() != tt.expected.Error() {
			t.Errorf("RegisterEventProcessor(%v): expected %s, actual %s", tt.given, tt.expected.Error(), actual)
		}
	}
}
