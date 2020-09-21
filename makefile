SHELL := /bin/bash

build:
	go build -o ./bin/httpservertest

run: build
	@echo starting server
	@GODEBUG=http2debug=2 ./bin/httpservertest > server.log 2>&1 &

stop:
	@echo stopping server
	kill -s INT $$(ps -ax | grep httpservertest | grep -v grep | awk '{print $$1}') >/dev/null 2>&1 || true

curl:
	@echo sending POST to service
	./curl-post.sh