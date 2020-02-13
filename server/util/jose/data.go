package jose

import "github.com/yixy/gateway/server/util/response"

const (
	FORMAT       = ""
	CHARSET      = "utf8"
	ENCRYPT_TYPE = ""
	HASH_TYPE    = ""
	HASHED       = ""
)

type Data struct {
	Format      string             `json:"format"`
	Charset     string             `json:"charset"`
	EncryptType string             `json:"encrypt_type"`
	HashType    string             `json:"hash_type"`
	Hashed      string             `json:"hashed"`
	Resp        *response.Response `json:"resp,omitempty"`
}
