package handler

import (
	"net/http"

	"github.com/yixy/gateway/server/util/response"

	"github.com/yixy/gateway/server/util/jose"

	"github.com/yixy/gateway/cfg"

	"github.com/SermoDigital/jose/jws"
	uuid "github.com/satori/go.uuid"
	//"github.com/yixy/tianmu-security/encrypt"
)

func APIhandler(w http.ResponseWriter, req *http.Request) {
	uuid := uuid.NewV4()
	data := &jose.Data{
		Format:      jose.FORMAT,
		Charset:     jose.CHARSET,
		EncryptType: jose.ENCRYPT_TYPE,
		HashType:    jose.HASH_TYPE,
		Hashed:      jose.HASHED,
		Resp:        nil,
	}
	j := req.Header.Get("Gw-Jwt")
	jwt, err := jws.ParseJWT([]byte(j))
	if err != nil {
		response.Return403Err(uuid, w, data, response.INNER_ERROR, err)
		return
	}
	//c := w.Claims()
	//fmt.Println(c)
	signingMethod := jws.GetSigningMethod("RS512")
	if err := jwt.Validate(cfg.PubKey, signingMethod); err != nil {
		panic(err)
	}

	//sign, err := encrypt.RsaSign(orignData, priKey, algorithm)
	//encrypt.RsaVerify(originData, cfg.PubKey, sign, algorithm)
	//req.Header.Get("response_biz_content")
	//req.Header.Get("sign")

	//req.Header.Get("return_code")
	//req.Header.Get("return_msg")
	//req.Header.Get("msg_id")

	//_, err := w.Write([]byte("404"))
	//if err != nil {
	//}
}
