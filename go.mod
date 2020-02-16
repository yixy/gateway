module github.com/yixy/gateway

go 1.13

require (
	github.com/SermoDigital/jose v0.9.2-0.20180104203859-803625baeddc
	github.com/mitchellh/mapstructure v1.1.2
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
	go.uber.org/zap v1.10.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/SermoDigital/jose v0.9.2-0.20180104203859-803625baeddc => github.com/yixy/jose v1.1.0
