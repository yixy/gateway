package resp

import (
	"encoding/json"
	"net/http"

	"github.com/yixy/gateway/log"
	"go.uber.org/zap"
)

type Response struct {
	ReturnCode ReturnCode  `json:"return_code"`
	ReturnMsg  string      `json:"return_msg"`
	Content    interface{} `json:"content,omitempty"`
}

type ReturnCode int

const (
	OK                             = 0
	FORBIDDEN_JWT_INVALID          = -40301
	FORBIDDEN_JWT_ALG_404          = -40302
	FORBIDDEN_JWT_ALG_ERR          = -40303
	FORBIDDEN_JWT_SUB_404          = -40304
	FORBIDDEN_JWT_SUB_ERR          = -40305
	FORBIDDEN_JWT_AUD_404          = -40306
	FORBIDDEN_JWT_AUD_ERR          = -40307
	FORBIDDEN_JWT_IAT_404          = -40308
	FORBIDDEN_JWT_IAT_ERR          = -40309
	FORBIDDEN_JWT_JWI_404          = -40310
	FORBIDDEN_JWT_JWI_ERR          = -40311
	FORBIDDEN_JWT_DATA_ERR         = -40312
	FORBIDDEN_JWT_HASHTYPE         = -40313
	FORBIDDEN_JWT_FORMAT           = -40314
	FORBIDDEN_JWT_VERRIFY          = -40315
	INTERNAL_SERVER_ERROR_JWS      = -50001
	INTERNAL_SERVER_ERROR_PROXYREQ = -50002
	BAD_GATEWAY_CONN               = -50201
	BAD_GATEWAY_CONN2              = -50202
	BAD_GATEWAY_IO_ERR             = -50203
	GATEWAY_TIMEOUT                = -50401
)

var Resp = map[ReturnCode]string{
	OK:                             "OK",
	FORBIDDEN_JWT_INVALID:          "FORBIDDEN: Parse JWT from header(Gateway-Jwt) error",
	FORBIDDEN_JWT_ALG_404:          "FORBIDDEN: JWT.header().alg not found",
	FORBIDDEN_JWT_ALG_ERR:          "FORBIDDEN: JWT.header().alg is invalid",
	FORBIDDEN_JWT_SUB_404:          "FORBIDDEN: JWT.payload().sub not found",
	FORBIDDEN_JWT_SUB_ERR:          "FORBIDDEN: JWT.payload().sub is not matched requestURI",
	FORBIDDEN_JWT_AUD_404:          "FORBIDDEN: JWT.payload().aud not found",
	FORBIDDEN_JWT_AUD_ERR:          "FORBIDDEN: JWT.payload().aud is not matched icbc",
	FORBIDDEN_JWT_IAT_404:          "FORBIDDEN: JWT.payload().iat not found",
	FORBIDDEN_JWT_IAT_ERR:          "FORBIDDEN: JWT.payload().iat is too old",
	FORBIDDEN_JWT_JWI_404:          "FORBIDDEN: JWT.payload().jwi not found",
	FORBIDDEN_JWT_JWI_ERR:          "FORBIDDEN: JWT.payload().jwi is empty",
	FORBIDDEN_JWT_DATA_ERR:         "FORBIDDEN: JWT.payload().data is invalid",
	FORBIDDEN_JWT_HASHTYPE:         "FORBIDDEN: JWT.payload().data.hash_type is invalid",
	FORBIDDEN_JWT_FORMAT:           "FORBIDDEN: JWT.payload().data.format is invalid",
	FORBIDDEN_JWT_VERRIFY:          "FORBIDDEN: JWT is expired or signing verify failed",
	INTERNAL_SERVER_ERROR_JWS:      "INTERNAL SERVER ERROR: Parse JWS error",
	INTERNAL_SERVER_ERROR_PROXYREQ: "INTERNAL SERVER ERROR: new proxy request error",
	BAD_GATEWAY_CONN:               "BAD_GATEWAY: Send request ro upstream server error",
	BAD_GATEWAY_CONN2:              "BAD_GATEWAY: upstream server response is nil",
	BAD_GATEWAY_IO_ERR:             "BAD_GATEWAY: copy upstream server response err",
	GATEWAY_TIMEOUT:                "GATEWAY TIMEOUT: Call upstream server timeout",
}

func GetResponse(appId string, urlPath string, jwtId string, signAlg string,
	data *Data, returnCode ReturnCode, err error, zuuid zap.Field) (jwt string, body []byte) {
	returnMsg := Resp[returnCode]
	log.Logger.Error(returnMsg, zap.Any("return_code", returnCode), zap.Error(err))
	//set data
	data.Resp = &Response{returnCode, returnMsg, nil}
	body, err = json.Marshal(data.Resp)
	if err != nil {
		log.Logger.Error("json.Marshal data.Resp error", zap.Error(err), zuuid)
	} else {
		data.Hashed, err = ShaSumS(body, data.HashType)
		if err != nil {
			log.Logger.Error("ShaSumS() calculate err", zap.Error(err), zuuid)
		}
	}
	//generate jwt
	jwt, err = GetRespJwt(appId, urlPath, jwtId, signAlg, data)
	if err != nil {
		log.Logger.Error("GetRespJwt error", zap.Error(err), zuuid)
	}
	log.Logger.Error("jwt", zap.String(GATEWAY_JWT, jwt), zuuid)
	return jwt, body
}

//403 Forbidden
func Return403Err(
	w http.ResponseWriter, appId string, urlPath string, jwtId string,
	signAlg string, data *Data, returnCode ReturnCode, err error, zuuid zap.Field) {
	jwt, body := GetResponse(appId, urlPath, jwtId, signAlg, data, returnCode, err, zuuid)
	w.Header().Set(GATEWAY_JWT, jwt)
	w.WriteHeader(403)
	w.Write(body)
}

//500 Internal Server Error
func Return500Err(
	w http.ResponseWriter, appId string, urlPath string, jwtId string,
	signAlg string, data *Data, returnCode ReturnCode, err error, zuuid zap.Field) {
	jwt, body := GetResponse(appId, urlPath, jwtId, signAlg, data, returnCode, err, zuuid)
	w.Header().Set(GATEWAY_JWT, jwt)
	w.WriteHeader(500)
	w.Write(body)
}

//502 Bad Gateway
func Return502Err(
	w http.ResponseWriter, appId string, urlPath string, jwtId string,
	signAlg string, data *Data, returnCode ReturnCode, err error, zuuid zap.Field) {
	jwt, body := GetResponse(appId, urlPath, jwtId, signAlg, data, returnCode, err, zuuid)
	w.Header().Set(GATEWAY_JWT, jwt)
	w.WriteHeader(502)
	w.Write(body)
}

//504 Gateway Timeout
func Return504Err(
	w http.ResponseWriter, appId string, urlPath string, jwtId string,
	signAlg string, data *Data, returnCode ReturnCode, err error, zuuid zap.Field) {
	jwt, body := GetResponse(appId, urlPath, jwtId, signAlg, data, returnCode, err, zuuid)
	w.Header().Set(GATEWAY_JWT, jwt)
	w.WriteHeader(504)
	w.Write(body)
}
