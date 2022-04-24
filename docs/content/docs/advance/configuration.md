---
weight: 1
---

# Configuration

配置为两份文件分别为 service config 和 registry config

- service config 关注服务启停以及声明周期中需要的各类配置
- registry config 关注服务注册相关配置

## server config definition

```proto
// Service 为一个服务的全部配置
message Service {
    string name = 1;
    string version = 2;
    optional Server http = 3;
    optional Server grpc = 4;
    Log log = 5;
    map<string, string> metadata = 6;
    Redis redis = 7;
    MQ mq = 8;
}

message Server {
    string scheme = 1;
    string addr = 2;
    int32 port = 3;
}


enum Level {
    DEBUG = 0;
    INFO = 1;
    WARING = 2;
    ERROR = 3;
    FATAL = 4;
}

message Log {
    optional string log_path = 1;
    repeated Level level = 2;
}

message Redis {
    string addr = 1;
    string password = 2;
    int32 max_conns = 3;
    int32 min_idle_conns = 4;
    google.protobuf.Duration dial_timeout = 5;
    google.protobuf.Duration idle_timeout = 6;
}

message MQ {
    repeated string addr = 1;
    int32 max_retry = 2;
}
```

## registry config definition

```proto
message RegistryInfo {
    repeated string addr = 1;
    string scheme = 2;
    google.protobuf.Duration dial_timeout_sec = 3;
    google.protobuf.Duration dial_keep_alive_time_sec = 4;
    google.protobuf.Duration dial_keep_alive_timeout_sec = 5;
}

message Registry {
    string name = 1;
    oneof reg {
        RegistryInfo consul = 2;
        RegistryInfo etcd = 3;
    }
}
```
