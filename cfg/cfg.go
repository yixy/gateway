package cfg

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

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
	JWT_TIMEOUT      = "JWT_TIMEOUT"
	CLIENT_TIMEOUT   = "CLIENT_TIMEOUT"
)

const (
	TIMEOUT       = 20  //server read write timeout
	CTIMEOUT      = 15  //client timeout
	TIMEOUT_VALID = 300 //jwt check timeout
)

//execute binary path
var Dir string

var (
	Pid           string
	LogFile       string
	PriKey        interface{}
	PubKey        interface{}
	Port          string
	Rtimeout      int64
	Wtimeout      int64
	ClientTimeout int64
	ShutTimeout   int64
	JwtTimeout    int64
)

func ReadCfg() error {
	var err error
	Pid := strconv.Itoa(os.Getpid())
	fmt.Println("Pid", Pid)
	fmt.Println("========= print config file =========")
	for _, key := range viper.AllKeys() {
		fmt.Println(key, viper.Get(key))
	}
	fmt.Println("=========       end         =========")
	LogFile = viper.GetString(LOG_FILE)
	priFile := viper.GetString(PRI_KEY_FILE)
	pubFile := viper.GetString(PUB_KEY_FILE)
	Port = viper.GetString(PORT)
	Rtimeout = viper.GetInt64(READ_TIMEOUT)
	Wtimeout = viper.GetInt64(WRITE_TIMEOUT)
	ShutTimeout = viper.GetInt64(SHUTDOWN_TIMEOUT)
	JwtTimeout = viper.GetInt64(JWT_TIMEOUT)
	ClientTimeout = viper.GetInt64(CLIENT_TIMEOUT)
	if isEmpty(LogFile, priFile, pubFile, Port) {
		return errors.New("LogFile, Port must not be empty.")
	}
	if !filepath.IsAbs(LogFile) {
		LogFile = filepath.Join(Dir, LogFile)
	}
	if !filepath.IsAbs(priFile) {
		priFile = filepath.Join(Dir, priFile)
	}
	if !filepath.IsAbs(pubFile) {
		pubFile = filepath.Join(Dir, pubFile)
	}
	PriKey, err = LoadKey(priFile)
	if err != nil {
		return err
	}
	PubKey, err = LoadKey(pubFile)
	if err != nil {
		return err
	}
	if Rtimeout <= 0 {
		Rtimeout = TIMEOUT
	}
	if Wtimeout <= 0 {
		Wtimeout = TIMEOUT
	}
	if ShutTimeout <= 0 {
		ShutTimeout = TIMEOUT * 2
	}
	if JwtTimeout <= 0 {
		JwtTimeout = TIMEOUT_VALID
	}
	if ClientTimeout <= 0 {
		ClientTimeout = CTIMEOUT
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

func LoadKey(file string) (interface{}, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	} else if block.Type == "PUBLIC KEY" {
		//公钥标准：x.509
		return x509.ParsePKIXPublicKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		//私钥标准：pkcs #8
		return x509.ParsePKCS8PrivateKey(block.Bytes)
	} else {
		return nil, errors.New("PEM format error, must be pkcs#8 or x.509")
	}
}
