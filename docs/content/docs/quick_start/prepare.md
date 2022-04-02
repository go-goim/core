---
weight: 2
title: "Prepare"
---

## requirement

### environment

- apache/rocketmq 4.6.0+
- consul 1.11.4+
- redis 2.0+

> 关于部署 rocketmq：docker 部署 rocketmq 过程中遇到过一些问题，如果你有疑问可以参考这篇文章 [Docker 部署 RocketMQ](https://yusank.space/posts/rocketmq-deploy/)

### config

msg service 为例，`apps/msg/config/config.yaml`:

    name: goim.msg.service
    version: v0.0.0
    grpc:
    scheme: grpc
    port: 18063
    log:
    level:
    - INFO
    - DEBUG
    metadata:
    grpcSrv: yes
    redis:
    addr: 127.0.0.1:6379
    mq:
    addr:
        - 127.0.0.1:9876

`apps/msg/config/registry.yaml` :

    consul:
    addr:
        - 127.0.0.1:8500
    scheme: http

根据自己的环境去修改各个组件的地址和端口。
