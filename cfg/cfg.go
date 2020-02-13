package cfg

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	LOG_FILE         = "LOG_FILE"
	PRI_KEY_FILE     = "PRI_KEY_FILE"
	PUB_KEY_FILE     = "PUB_KEY_FILE"
	PORT             = "PORT"
	READ_TIMEOUT     = "READ_TIMEOUT"
	WRITE_TIMEOUT    = "WRITE_TIMEOUT"
	SHUTDOWN_TIMEOUT = "SHUTDOWN_TIMEOUT"
)

const TIMEOUT = 20000

//execute binary path
var Dir string

var LogFile string
var PriFile string
var PubFile string
var PriKey []byte
var PubKey []byte
var Port string
var Rtimeout int64
var Wtimeout int64
var ShutTimeout int64

func CfgCheck() error {
	var err error
	fmt.Println("========= print config file =========")
	for _, key := range viper.AllKeys() {
		fmt.Println(key, viper.Get(key))
	}
	fmt.Println("=========       end         =========")
	LogFile = viper.GetString(LOG_FILE)
	PriFile = viper.GetString(PRI_KEY_FILE)
	PubFile = viper.GetString(PUB_KEY_FILE)
	Port = viper.GetString(PORT)
	Rtimeout = viper.GetInt64(READ_TIMEOUT)
	Wtimeout = viper.GetInt64(WRITE_TIMEOUT)
	ShutTimeout = viper.GetInt64(SHUTDOWN_TIMEOUT)
	if isEmpty(LogFile, PriFile, PubFile, Port) {
		return errors.New("LogFile, Port must not be empty.")
	}
	if !filepath.IsAbs(LogFile) {
		LogFile = filepath.Join(Dir, LogFile)
	}
	if !filepath.IsAbs(PriFile) {
		PriFile = filepath.Join(Dir, PriFile)
	}
	if !filepath.IsAbs(PubFile) {
		PubFile = filepath.Join(Dir, PubFile)
	}
	PriKey, err = LoadKey(PriFile)
	if err != nil {
		return err
	}
	PubKey, err = LoadKey(PubFile)
	if err != nil {
		return err
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

func LoadKey(file string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	} else if block.Type != "PUBLIC KEY" && block.Type != "PRIVATE KEY" {
		//私钥标准：pkcs #8
		//公钥标准：x.509
		return nil, errors.New("failed to decode PEM block containing public/private key")
	}
	return block.Bytes, nil
}
