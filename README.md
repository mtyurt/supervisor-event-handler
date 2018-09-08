# supervisor-event-handler

This small library abstracts away supervisor's event handling protocol and provides
an easy way to process events. For more detailed information about events, please check [here](http://supervisord.org/events.html).

# features

- support generic events, like `PROCESS_STATE` to handle `PROCESS_STATE*` events
- run processors via goroutines to avoid buffer overflow as much as possible

# installation

```
go get -u github.com/mtyurt/supervisor-event-handler
```

# usage

The `example/` directory contains a fully working application with Docker.


```go
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
```

- `handler.Start()` is a blocking process, it will run until the process is killed.
- Event handlers are run via goroutines to not overflow the supervisor event buffer. If the process cannot keep up with incoming events, the oldest event will be discarded by supervisor.
- Supervisor's event handler mechanism requires the process to print to stdout, so do not print to stdout.

# licence

The BSD 3-Clause License - see [LICENSE](https://github.com/mtyurt/supervisor-event-handler/blob/master/LICENSE) for more details.
