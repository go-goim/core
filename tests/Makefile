SHELL:=/usr/bin/env bash

ADDR ?= 127.0.0.1:18071
UID ?= user1

.PHONEY: run-gui
run-gui:
	go run client-gui/*.go -s $(ADDR) -u $(UID)

.PHONEY: run-bench
run-bench:
	go run client-bench/*.go