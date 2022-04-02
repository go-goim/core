---
weight: 1
---

# Quick Start

## run

```shell
# run msg service
$ make run Srv=msg
# run gateway service
$ make run Srv=gateway
# run push service
$ make run Srv=push
```

## other make command

```shell
make help

Usage:
  make <target>

Development
  vet              Run go vet against code.
  lint             Run go lint against code.
  test             Run test against code.

Generate
  protoc           Run protoc command to generate pb code.

Build
  build            build provided server
  build-all        build all apps

Docker
  docker-build     build docker image

Run
  run              run provided server

General
  help             Display this help.
```
