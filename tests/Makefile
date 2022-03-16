SHELL:=/usr/bin/env bash

ADDR ?= 127.0.0.1:8071
UID ?= user1
TOUID ?= user2

.PHONEY: run
run:
	go run client/*.go -s $(ADDR) -u $(UID) -t $(TOUID)