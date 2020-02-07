package cfg

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	LOG_FILE         = "LOG_FILE"
	PORT             = "PORT"
	READ_TIMEOUT     = "READ_TIMEOUT"
	WRITE_TIMEOUT    = "WRITE_TIMEOUT"
	SHUTDOWN_TIMEOUT = "SHUTDOWN_TIMEOUT"
)

const TIMEOUT = 20000

//execute binary path
var Dir string

var LogFile string
var Port string
var Rtimeout int64
var Wtimeout int64
var ShutTimeout int64

func CfgCheck() error {
	fmt.Println("========= print config file =========")
	for _, key := range viper.AllKeys() {
		fmt.Println(key, viper.Get(key))
	}
	fmt.Println("=========       end         =========")
	LogFile = viper.GetString(LOG_FILE)
	Port = viper.GetString(PORT)
	Rtimeout = viper.GetInt64(READ_TIMEOUT)
	Wtimeout = viper.GetInt64(WRITE_TIMEOUT)
	ShutTimeout = viper.GetInt64(SHUTDOWN_TIMEOUT)
	if isEmpty(LogFile, Port) {
		return errors.New("LogFile, Port must not be empty.")
	}
	if Rtimeout == 0 {
		Rtimeout = TIMEOUT
	}
	if Wtimeout == 0 {
		Wtimeout = TIMEOUT
	}
	if ShutTimeout == 0 {
		ShutTimeout = TIMEOUT * 2
	}
	return nil
}

func isEmpty(keys ...string) (result bool) {
	result = false
	for _, key := range keys {
		if key == "" {
			result = true
		}
	}
	return result
}
