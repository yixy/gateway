//test program for generating datas.
package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/yixy/gateway/cfg"
	"github.com/yixy/gateway/server/util/jose"

	"github.com/SermoDigital/jose/jws"
)

func main() {
	//body hashed
	hashed, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	pubKey, err := cfg.LoadKey("pub.pem")
	if err != nil {
		panic(err)
	}

	priKey, err := cfg.LoadKey("pri.pem")
	if err != nil {
		panic(err)
	}

	rsaPub, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		panic(err)
	}

	rsaPriv, err := x509.ParsePKCS8PrivateKey(priKey)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	claims := jws.Claims{}
	claims.SetIssuer("appid00001") //appid
	claims.SetIssuedAt(now)        //timestamp
	claims.SetNotBefore(now)
	claims.SetExpiration(now.Add(time.Duration(24) * time.Hour))
	claims.SetSubject("/api/xxx/yyy/v2")
	claims.SetAudience("icbc")
	claims.SetJWTID("msgid-abcdefghijklmnopqrstuvwxyz") //msgid
	//
	data := jose.Data{
		Charset:     jose.CHARSET,
		Format:      jose.FORMAT,
		EncryptType: jose.ENCRYPT_TYPE,
		Hashed:      string(hashed),
		HashType:    os.Args[1],
		Resp:        nil,
	}

	claims.Set("data", data)

	signingMethod := jws.GetSigningMethod("RS512")
	j := jws.NewJWT(claims, signingMethod)
	b, err := j.Serialize(rsaPriv)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", string(b))

	//time.Sleep(time.Second * 2)

	w, err := jws.ParseJWT(b)
	if err != nil {
		panic(err)
	}
	if err := w.Validate(rsaPub, signingMethod); err != nil {
		panic(err)
	}
}
