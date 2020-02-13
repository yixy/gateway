package jose

import (
	"time"

	"github.com/SermoDigital/jose/jws"
)

type Conf struct {
	Method string // 加密算法
	Key    string // 加密key
	Issuer string // 签发者
	Expire int64  // 签名有效期
}

var conf = Conf{
	Method: "RS512",
	Key:    "sahjdjsgaudsiudhuywge",
	Issuer: "ICBC-API",
	Expire: 100,
}

// GetJWT 获取json web token
func GetJWT(data map[string]interface{}) (token string, err error) {
	payload := jws.Claims{}
	for k, v := range data {
		payload.Set(k, v)
	}
	now := time.Now()
	payload.SetIssuer(conf.Issuer)
	payload.SetExpiration(now.Add(time.Duration(conf.Expire) * time.Second))
	payload.SetSubject("/api/test") //api-url
	payload.SetAudience("icbc")     //icbc
	payload.SetNotBefore(now)
	payload.SetIssuedAt(now)
	payload.SetJWTID("aaaaaaaa") //msgid
	signingMethod := jws.GetSigningMethod(conf.Method)
	signingMethod.Sign()
	jwtObj := jws.NewJWT(payload, signingMethod)
	tokenBytes, err := jwtObj.Serialize([]byte(conf.Key))
	if err != nil {
		return
	}
	token = string(tokenBytes)
	return
}

// VerifyJWT 验证json web token
func VerifyJWT(token string) (ret bool, err error) {
	jwtObj, err := jws.ParseJWT([]byte(token))
	if err != nil {
		return
	}
	err = jwtObj.Validate([]byte(conf.Key), jws.GetSigningMethod(conf.Method))
	if err == nil {
		ret = true
	}
	return
}
