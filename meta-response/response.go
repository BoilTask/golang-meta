package metaresponse

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"meta/error-code"
	metaerror "meta/meta-error"
	metaformat "meta/meta-format"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

func NewResponse[T metaerrorcode.Numeric](ctx *gin.Context, code T, data ...interface{}) {
	response := Response{
		Code: int(code),
	}
	if data != nil {
		if len(data) > 1 {
			response.Data = data
		} else if len(data) > 0 {
			response.Data = data[0]
		}
	}
	ctx.JSON(http.StatusOK, response)
}

func NewResponseError(ctx *gin.Context, err error, data ...interface{}) {
	if err == nil {
		NewResponse(ctx, metaerrorcode.Success, data...)
		return
	}
	_ = ctx.Error(err)
	response := Response{
		Code: metaerror.GetErrorCodeFromError(err),
	}
	if data != nil {
		if len(data) > 1 {
			response.Data = data
		} else if len(data) > 0 {
			response.Data = data[0]
		}
	}
	ctx.JSON(http.StatusOK, response)
}

func NewResponseWrapCode[T metaerrorcode.Numeric](
	ctx *gin.Context,
	err error,
	code T,
	data interface{},
	format ...any,
) {
	if err == nil {
		slog.Info("NewResponseWrapCode", "message", metaformat.Format(format...))
		NewResponse(ctx, code, data)
		return
	}
	NewResponseError(
		ctx, metaerror.WrapCode(
			err,
			code,
			metaformat.Format(format...),
		), data,
	)
}
