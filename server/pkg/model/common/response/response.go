package response

import (
	"fmt"
	"net/http"
	"server/pkg/global"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

const SUCCESS = 20000

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func Result(code int, data any, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *gin.Context, err error, data any, msg ...string) {
	if err != nil {
		global.MPA_LOG.Error(fmt.Errorf("%-v", err).Error())
		coder := errors.ParseCoder(err)

		c.JSON(http.StatusOK, ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})

		return
	}

	if len(msg) > 0 {
		c.JSON(http.StatusOK, Response{SUCCESS, data, msg[0]})
	} else {
		c.JSON(http.StatusOK, Response{Code: SUCCESS, Data: data})
	}
}

func SuccessResponse(c *gin.Context, data any, msg ...string) {
	if len(msg) > 0 {
		c.JSON(http.StatusOK, Response{Code: SUCCESS, Data: data, Message: msg[0]})
	} else {
		c.JSON(http.StatusOK, Response{Code: SUCCESS, Data: data, Message: "SUCCESS"})
	}
}

func FailResponse(c *gin.Context, err error) {
	coder := errors.ParseCoder(err)

	// 非200状态码，记录日志
	if coder.HTTPStatus() != 200 {
		global.MPA_LOG.Error(fmt.Errorf("%-v", err).Error())
	}

	c.JSON(http.StatusOK, ErrResponse{
		Code:      coder.Code(),
		Message:   coder.String(),
		Reference: coder.Reference(),
	})
}
