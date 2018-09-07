package supervisorevent

import "testing"

func TestParseHeader(t *testing.T) {
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
		actual := parseHeaderTokens(tt.given)
		if actual != tt.expected {
			t.Errorf("parseHeaderTokens(%s): expected %v, actual %v", tt.given, tt.expected, actual)
		}
	}
}
