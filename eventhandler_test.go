package eventhandler

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseHeader(t *testing.T) {
	h := EventHandler{}
	var tests = []struct {
		name     string
		expected HeaderTokens
		given    string
	}{
		{
			"usual_order_header",
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
			"reversed_random_order_header",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := h.parseHeaderTokens(tt.given)
			if actual != tt.expected {
				t.Errorf("h.parseHeaderTokens(%s): expected %v, actual %v", tt.given, tt.expected, actual)
			}

		})
	}
}

func TestHandleEvent(t *testing.T) {
	h := EventHandler{}
	var tests = []struct {
		name     string
		expected error
		given    struct {
			event   string
			handler EventProcessor
		}
	}{
		{"state_change_event", nil, struct {
			event   string
			handler EventProcessor
		}{
			"PROCESS_STATE",
			nil,
		}},
		{"invalid_event_name", fmt.Errorf("invalidevent is not a valid event! Valid events are: %v", ValidEventNames), struct {
			event   string
			handler EventProcessor
		}{
			"invalidevent",
			nil,
		}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := h.HandleEvent(tt.given.event, tt.given.handler)
			if (actual != nil && tt.expected != nil) && actual.Error() != tt.expected.Error() {
				t.Errorf("HandleEvent(%v): expected %s, actual %s", tt.given, tt.expected.Error(), actual)
			}

		})
	}
}

type readResponse struct {
	header  HeaderTokens
	payload map[string]string
	err     error
}

func TestReadHeaderAndPayload(t *testing.T) {
	h := EventHandler{}

	var tests = []struct {
		name     string
		expected readResponse
		given    string
	}{
		{
			"running_state_test",
			readResponse{
				HeaderTokens{"3.0", "supervisor", "21", "listener", "10", "PROCESS_STATE_RUNNING", 58},
				map[string]string{
					"processname": "cat",
					"groupname":   "cat",
					"from_state":  "STARTING",
					"pid":         "2766",
				},
				nil,
			},
			"ver:3.0 server:supervisor serial:21 pool:listener poolserial:10 eventname:PROCESS_STATE_RUNNING len:58\nprocessname:cat groupname:cat from_state:STARTING pid:2766",
		},
		{
			"no_payload",
			readResponse{
				HeaderTokens{"3.0", "supervisor", "21", "listener", "10", "SUPERVISOR_STATE_CHANGE_RUNNING", 0},
				map[string]string{},
				nil,
			},
			"ver:3.0 server:supervisor serial:21 pool:listener poolserial:10 eventname:SUPERVISOR_STATE_CHANGE_RUNNING len:0\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			header, payload, err := h.readHeaderAndPayload(bufio.NewReader(strings.NewReader(tt.given)))
			if header != tt.expected.header || !cmp.Equal(payload, tt.expected.payload) || err != tt.expected.err {
				t.Errorf("h.readHeaderAndPayload(%v): expected %v, actual %v", tt.given, tt.expected, readResponse{header, payload, err})
			}
		})
	}
}

func TestProcessEvent(t *testing.T) {
	actual := ""
	h := EventHandler{}
	if err := h.HandleEvent("PROCESS_STATE", func(header HeaderTokens, payload map[string]string) {
		actual = header.EventName + payload["processname"]
	}); err != nil {
		t.Errorf(err.Error())
	}

	if err := h.HandleEvent("PROCESS_GROUP_ADDED", func(header HeaderTokens, payload map[string]string) {
		actual = header.EventName + payload["groupname"]
	}); err != nil {
		t.Errorf(err.Error())
	}
	var tests = []struct {
		expected     string
		givenHeader  HeaderTokens
		givenPayload map[string]string
	}{
		{
			"PROCESS_STATE_RUNNINGcat",
			HeaderTokens{EventName: "PROCESS_STATE_RUNNING"},
			map[string]string{"processname": "cat"},
		},
		{
			"PROCESS_GROUP_ADDEDbat",
			HeaderTokens{EventName: "PROCESS_GROUP_ADDED"},
			map[string]string{"groupname": "bat"},
		},
		{
			"",
			HeaderTokens{EventName: "PROCESS_LOG_STDOUT"},
			map[string]string{"data": "abc"},
		},
	}
	for _, tt := range tests {
		tt := tt
		actual = ""
		h.processEvent(tt.givenHeader, tt.givenPayload)
		if actual != tt.expected {
			t.Errorf("h.process(%v, %v): expected %s, actual %s", tt.givenHeader, tt.givenPayload, tt.expected, actual)
		}
	}

}
