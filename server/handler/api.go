package handler

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/SermoDigital/jose/jws"
	uuid "github.com/satori/go.uuid"
	"github.com/yixy/gateway/cfg"
	"github.com/yixy/gateway/log"
	"github.com/yixy/gateway/server/util/resp"
	"go.uber.org/zap"
	//"github.com/yixy/tianmu-security/encrypt"
)

const MAX_CONNS_PER_HOST = 100

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: MAX_CONNS_PER_HOST,
	},
	Timeout: time.Duration(cfg.CTIMEOUT) * time.Millisecond,
}

func ServiceHandler(w http.ResponseWriter, req *http.Request) {
	var (
		//request jwt header properties
		alg string = resp.RS512
		//request jwt payload properties
		iss  string
		sub  string //requestURI path[?query]
		aud  string
		iat  time.Time
		jti  string     //msgid
		data *resp.Data = &resp.Data{
			Format:      resp.FORMAT,
			Charset:     resp.CHARSET,
			EncryptType: resp.ENCRYPT_TYPE,
			HashType:    resp.SHA512,
			Hashed:      resp.HASHED,
			Resp:        nil,
		}
	)

	now := time.Now()

	uuid := uuid.NewV4().String()
	zuuid := zap.String("uuid", uuid)
	requestURI := req.RequestURI
	log.Logger.Info("ServiceHandler start.", zuuid, zap.String("RequestURI", requestURI))
	for i, v := range req.Header {
		log.Logger.Info("request header", zuuid, zap.Strings(i, v))
	}

	//get jwt
	token := req.Header.Get(resp.GATEWAY_JWT)
	jwt, err := jws.ParseJWT([]byte(token))
	if err != nil {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_INVALID, err, zuuid)
		return
	}

	//check header.alg
	j, ok := jwt.(jws.JWS)
	if !ok {
		resp.Return500Err(w, iss, requestURI, jti, alg, data, resp.INTERNAL_SERVER_ERROR_JWS, nil, zuuid)
		return
	}
	alg, ok = j.Protected().Get(resp.ALG).(string)
	if !ok {
		alg = resp.RS512
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_ALG_404, nil, zuuid)
		return
	}
	if !resp.CheckAlg(alg) {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_ALG_ERR, nil, zuuid)
		return
	}

	//verify payload.issuer : appid
	iss, ok = jwt.Claims().Issuer()
	if !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_ALG_ERR, nil, zuuid)
		return
	}
	//TODO
	//check appId
	if ok, returnCode := checkAppId(iss, alg, data); !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, returnCode, nil, zuuid)
		return
	}

	//verify payload.sub : requestURI
	sub, ok = jwt.Claims().Subject()
	if !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_SUB_404, nil, zuuid)
		return
	}
	if sub != requestURI {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_SUB_ERR, nil, zuuid)
		return
	}

	//verify payload.aud
	auds, ok := jwt.Claims().Audience()
	if !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_AUD_404, nil, zuuid)
		return
	}
	aud = auds[0]
	if strings.ToLower(aud) != resp.ICBC {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_AUD_ERR, nil, zuuid)
		return
	}

	//verify payload.iat
	iat, ok = jwt.Claims().IssuedAt()
	if !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_IAT_404, nil, zuuid)
		return
	}
	if now.After(iat.Add(time.Second * time.Duration(cfg.JwtTimeout))) {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_IAT_ERR, nil, zuuid)
		return
	}

	//verify payload.jti
	jti, ok = jwt.Claims().JWTID()
	if !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_JWI_ERR, nil, zuuid)
		return
	}
	if jti == "" {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_JWI_ERR, nil, zuuid)
		return
	}

	//check payload.data
	err = mapstructure.Decode(jwt.Claims().Get(resp.DATA), data)
	if err != nil {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_DATA_ERR, nil, zuuid)
		return
	}
	log.Logger.Info("print payload().data", zap.Any("data", data), zuuid)
	if !resp.CheckShaType(data.HashType) {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_HASHTYPE, nil, zuuid)
		return
	}
	if data.Format != "stream" {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_FORMAT, nil, zuuid)
		return
	}

	//verify jwt's sign and jwt's expired
	signingMethod := jws.GetSigningMethod(alg)
	if err := jwt.Validate(cfg.PubKey, signingMethod); err != nil {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, resp.FORBIDDEN_JWT_VERRIFY, err, zuuid)
		return
	}

	proxyReq, err := http.NewRequest(req.Method, "", req.Body)
	if err != nil {

	}
	proxyResp, err := client.Do(proxyReq)
	if err != nil {

	}
	if proxyResp == nil {
		panic("")
	}

	//set headers

	_, err = io.Copy(w, proxyResp.Body)
	if err != nil {

	}

}

func checkAppId(appId, alg string, data *resp.Data) (bool, resp.ReturnCode) {
	//TODO get database data to check
	return true, resp.OK
}
