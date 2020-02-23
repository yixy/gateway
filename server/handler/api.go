package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/textproto"
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

var connection []string

const MAX_CONNS_PER_HOST = 100

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: MAX_CONNS_PER_HOST,
	},
	Timeout: time.Duration(cfg.ClientTimeout) * time.Second,
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
	if ok, returnCode, err := checkAppId(iss, alg, data); !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, returnCode, err, zuuid)
		return
	}

	//TODO
	//check API
	ok, returnCode, serviceUrl, err := checkApi(iss, req.URL.Path)
	if !ok {
		resp.Return403Err(w, iss, requestURI, jti, alg, data, returnCode, err, zuuid)
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

	proxyReq, err := http.NewRequest(req.Method, serviceUrl, req.Body)
	if err != nil {
		resp.Return500Err(w, iss, requestURI, jti, alg, data, resp.INTERNAL_SERVER_ERROR_PROXYREQ, err, zuuid)
		return
	}

	//set request headers
	connection = strings.Split(req.Header.Get("Connection"), ",")
	for k, v := range req.Header {
		//Hop-by-hop header, do not transfer
		if isHopByHop(k) {
			continue
		}
		flag := false
		for _, s := range connection {
			if k == textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(s)) {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		for _, vv := range v {
			proxyReq.Header.Add(k, vv)
		}
	}

	//send request to upstream server.
	log.Logger.Info("send request to upstream server.")
	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		t, ok := err.(interface {
			Timeout() bool
		})
		if ok && t.Timeout() {
			resp.Return504Err(w, iss, requestURI, jti, alg, data, resp.GATEWAY_TIMEOUT, err, zuuid)
			return
		}
		resp.Return502Err(w, iss, requestURI, jti, alg, data, resp.BAD_GATEWAY_CONN, err, zuuid)
		return
	}
	if proxyResp == nil {
		resp.Return502Err(w, iss, requestURI, jti, alg, data, resp.BAD_GATEWAY_CONN2, nil, zuuid)
		return
	}
	defer func() {
		err = proxyResp.Body.Close()
		if err != nil {
			log.Logger.Error("call defer", zap.Error(err))
		}
	}()

	//set response headers
	connection = strings.Split(proxyResp.Header.Get("Connection"), ",")
	respData := &resp.Data{}
	for k, v := range proxyResp.Header {
		//Hop-by-hop header, do not transfer
		if isHopByHop(k) {
			continue
		}
		if k == resp.GATEWAY_APIRESP {
			err := json.Unmarshal([]byte(v[0]), respData)
			if err != nil {
				resp.Return502Err(w, iss, requestURI, jti, alg, data, resp.BAD_GATEWAY_CONN2, err, zuuid)
				return
			}
			respData.HashType = data.HashType
			respData.Format = data.Format
			respData.Charset = data.Charset
			respData.EncryptType = data.EncryptType
			continue
		}
		flag := false
		for _, s := range connection {
			if k == textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(s)) {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	//set cookies
	//for _, value := range proxyResp.Request.Cookies() {
	//	w.Header().Add(value.Name, value.Value)
	//}

	respJwt, err := resp.GetRespJwt(iss, requestURI, jti, alg, respData)
	if err != nil {
		resp.Return500Err(w, iss, requestURI, jti, alg, data, resp.INTERNAL_SERVER_ERROR_GENJWT, err, zuuid)
		return
	}

	w.Header().Set(resp.GATEWAY_APIRESP, respJwt)
	w.WriteHeader(proxyResp.StatusCode)

	//read response from upstream server.
	log.Logger.Info("read response from upstream server.")
	_, err = io.Copy(w, proxyResp.Body)
	if err != nil {
		resp.Return502Err(w, iss, requestURI, jti, alg, data, resp.BAD_GATEWAY_IO_ERR, err, zuuid)
		return
	}

}

func checkAppId(appId, alg string, data *resp.Data) (ok bool, returnCode resp.ReturnCode, err error) {
	//TODO get database data to check
	return true, resp.OK, nil
}

func checkApi(appId, path string) (ok bool, returnCode resp.ReturnCode, serviceUrl string, err error) {
	//TODO
	return true, resp.OK, "http://10.168.0.6:7777/test", nil
}

func isHopByHop(header string) (result bool) {
	//Hop-by-hop header
	//* Connetion
	//* Keep-Alive
	//* Proxy-Authenticate
	//* Proxy-Authorization
	//* Trailer
	//* TE
	//* Transfer-Encoding
	//* Upgrade
	switch header {
	case "Connection":
		result = true
	case "Keep-Alive":
		result = true
	case "Proxy-Authenticate":
		result = true
	case "Proxy-Authorization":
		result = true
	case "Trailer":
		result = true
	case "TE":
		result = true
	case "Transfer-Encoding":
		result = true
	case "Upgrade":
		result = true
	default:
		result = false
	}
	return result
}
