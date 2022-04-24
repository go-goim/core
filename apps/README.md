# apps

Apps contains all executable services.

## services resources

resources:

| service/app | redis | mysql | hbase | mq | http(pprof Port+1) | grpc | tcp |
| :--- | :---: | :---: |:---: | :---: | :---: |:---: | :---: |
|gateway | Y |  N | N | Y | `:18071` | `:18073` | -  |
|msg | Y |  Y | N | Y | - | `:18063` | -  |
|push | Y |  N | N | N | `:18081` | `:18083` | -  |
|store | N |  N | Y | Y | - | - | -  |
|user | Y |  Y | N | N | `:18051` | `:18053` | -  |
