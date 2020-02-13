package response

import (
	"net/http"

	"github.com/yixy/gateway/log"
	"github.com/yixy/gateway/server/util/jose"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type Response struct {
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}

const (
	INNER_ERROR = -50041
)

var Resp = map[int]string{
	INNER_ERROR: "gateway inner error",
}

//403 Forbidden
func Return403Err(uuid uuid.UUID, w http.ResponseWriter, data *jose.Data, errType int, err error) {
	log.Logger.Error(Resp[errType], zap.Error(err))
}

//500 Internal Server Error
func Return500Err(uuid uuid.UUID, w http.ResponseWriter, data *jose.Data, errType int, err error) {
	log.Logger.Error(Resp[errType], zap.Error(err))

}

//502 Bad Gateway
func Return502Err(uuid uuid.UUID, w http.ResponseWriter, data *jose.Data, errType int, err error) {

}

//504 Gateway Timeout
func Return504Err(uuid uuid.UUID, w http.ResponseWriter, data *jose.Data, errType int, err error) {

}
