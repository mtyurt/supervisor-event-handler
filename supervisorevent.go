package supervisorevent

import (
	"log"
	"strconv"
	"strings"
)

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

// receives space separated {key}:{value} string pairs
// creates a map where key -> value
func parseTokensToMap(tokens string) map[string]string {
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
func parseHeaderTokens(header string) HeaderTokens {
	headerMap := parseTokensToMap(header)
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
