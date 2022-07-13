module github.com/go-goim/core

go 1.18

require (
	github.com/apache/rocketmq-client-go/v2 v2.1.0
	github.com/btcsuite/btcd/btcutil v1.1.1
	github.com/envoyproxy/protoc-gen-validate v0.6.7 // indirect
	github.com/gin-gonic/gin v1.8.1
	github.com/go-kratos/kratos/contrib/registry/etcd/v2 v2.0.0-20220301040457-03ad2b663624
	github.com/go-kratos/kratos/v2 v2.3.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/consul/api v1.13.0
	github.com/jroimartin/gocui v0.5.0
	github.com/mattn/go-colorable v0.1.12
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	go.etcd.io/etcd/client/v3 v3.5.4
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.3.3
	gorm.io/gorm v1.23.8
)

require (
	github.com/go-goim/api v0.0.5
	github.com/go-kratos/kratos/contrib/config/etcd/v2 v2.0.0-20220528114537-97c103a39562
	github.com/panjf2000/ants/v2 v2.5.0
)

require (
	github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.3.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-hclog v0.12.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/hashicorp/serf v0.9.6 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nsf/termbox-go v1.1.1 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/tidwall/gjson v1.2.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.etcd.io/etcd/api/v3 v3.5.4 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.4 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	stathat.com/c/consistent v1.0.0 // indirect
)

replace (
	github.com/json-iterator/go v1.1.11 => github.com/json-iterator/go v1.1.12
	github.com/tidwall/gjson v1.2.1 => github.com/tidwall/gjson v1.13.0
)
