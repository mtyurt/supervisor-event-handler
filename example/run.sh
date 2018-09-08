#!/bin/bash
set -ex

GOOS=linux GOARCH=amd64 vgo build -o example

docker build -f Dockerfile -t supervisor_event_handler_example .
docker run --rm -it \
    -v $(pwd):/root/eventhandler \
    supervisor_event_handler_example
