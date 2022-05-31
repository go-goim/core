#!/bin/bash

docker run -d --name rocketmq-dashboard --network docker-compose_default \
    -e "JAVA_OPTS=-Drocketmq.namesrv.addr=rmqnamesrv:9876" -p 8080:8080 \
    -t apacherocketmq/rocketmq-dashboard:latest