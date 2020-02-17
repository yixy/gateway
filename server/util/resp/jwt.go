package resp

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"time"

	"github.com/yixy/gateway/cfg"

	"github.com/SermoDigital/jose/jws"
)

const (
	FORMAT          = "stream"
	CHARSET         = "UTF-8"
	ENCRYPT_TYPE    = ""
	HASHED          = ""
	GATEWAY_JWT     = "Gateway-jwt"
	GATEWAY_APIRESP = "Gateway-Apiresp"
	ICBC            = "icbc"
	ALG             = "alg"
	DATA            = "data"
	RS256           = "RS256"
	RS384           = "RS384"
	RS512           = "RS512"
	SHA256          = "SHA256"
	SHA384          = "SHA384"
	SHA512          = "SHA512"
)

type Data struct {
	Format      string    `json:"format" mapstructure:"format"`
	Charset     string    `json:"charset" mapstructure:"charset"`
	EncryptType string    `json:"encrypt_type,omitempty" mapstructure:"encrypt_type"`
	HashType    string    `json:"hash_type" mapstructure:"hash_type"`
	Hashed      string    `json:"hashed" mapstructure:"hashed"`
	Resp        *Response `json:"resp,omitempty" mapstructure:"resp"`
}

func CheckAlg(alg string) bool {
	//check data from database
	//TODO
	switch alg {
	case RS256:
	case RS384:
	case RS512:
	default:
		return false
	}
	return true
}

func CheckShaType(sha string) bool {
	switch sha {
	case SHA256:
	case SHA384:
	case SHA512:
	default:
		return false
	}
	return true
}

func GetRespJwt(appId, urlPath, jwtId, signAlg string, data *Data) (jwt string, err error) {
	if !CheckAlg(signAlg) {
		return "", errors.New("signAlg invalid")
	}

	now := time.Now()
	claims := jws.Claims{}
	claims.SetIssuer(ICBC)  //appid
	claims.SetIssuedAt(now) //timestamp
	claims.SetNotBefore(now)
	claims.SetExpiration(now.Add(time.Duration(cfg.JwtTimeout) * time.Second))
	claims.SetSubject(urlPath) //apiurl path
	claims.SetAudience(appId)
	claims.SetJWTID(jwtId) //msgid

	claims.Set("data", *data)

	signingMethod := jws.GetSigningMethod(signAlg)
	j := jws.NewJWT(claims, signingMethod)
	b, err := j.Serialize(cfg.PriKey)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//for short data shasum
func ShaSumS(data []byte, shaType string) (string, error) {
	var h hash.Hash
	switch shaType {
	case SHA256:
		h = sha256.New()
	case SHA384:
		h = sha512.New384()
	case SHA512:
		h = sha512.New()
	default:
		return "", errors.New("shaType is invalid")
	}
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	//return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
