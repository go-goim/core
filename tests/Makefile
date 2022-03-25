SHELL:=/usr/bin/env bash

ADDR ?= 127.0.0.1:18071
UID ?= user1
TOUID ?= user2

.PHONEY: run-gui
run-gui:
	go run client-gui/*.go -s $(ADDR) -u $(UID) -t $(TOUID)

.PHONEY: run-bench
run-bench:
	go run client-bench/*.go