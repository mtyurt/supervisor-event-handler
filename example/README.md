# supervisor-event-handler-example
Here you can see a fully working example working on Docker. After satisfying requirements, please execute `run.sh`, which will to everything for you, and you can tail `eventhandler.log` file for event handler's output.

# requirements
- docker: https://docs.docker.com/install/
- vgo: `go get -u golang.org/x/vgo`

# supervisor configuration
`conf/supervisor.d/eventlistener.ini` contains a bare minimum configuration for the event listener process. Important points:

- Value of `events` can be multiple separated by a comma. Check event types [here](http://supervisord.org/events.html#event-types)
- If the process cannot keep up with incoming events, supervisor puts events to a buffer. If buffer is overflowed, the oldest event will be discarded when a new event comes.

# example output in log file
```
event: PROCESS_STATE_STARTING, payload: map[processname:simple_script groupname:simple_script from_state:EXITED tries:0]
event: PROCESS_STATE_RUNNING, payload: map[from_state:STARTING pid:26 processname:simple_script groupname:simple_script]
event: PROCESS_STATE_STOPPING, payload: map[processname:simple_script groupname:simple_script from_state:RUNNING pid:26]
```
