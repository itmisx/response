package response

import (
	"net/http"
	"strings"

	"github.com/itmisx/errorx"
	"github.com/itmisx/i18n"

	"github.com/gin-gonic/gin"
)

// JSON http json格式response响应
func JSON(ctx *gin.Context, v interface{}, e error) {
	lang := strings.Split(ctx.GetHeader("Accept-Language"), ",")[0]
	// 正常的返回
	if e == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "",
			"data": v,
		})
	}
	// 错误返回
	if e != nil {
		errorx, ok := e.(errorx.Error)
		// 自定义错误返回
		if ok {
			msg := i18n.T(lang, errorx.Code)
			if msg == "" {
				msg = e.Error()
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code": errorx.Code,
				"msg":  msg,
				"data": v,
			})
		} else {
			// 系统错误
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  e.Error(),
				"data": v,
			})
		}
	}
}
