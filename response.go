package response

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itmisx/errorx"
	"github.com/itmisx/i18n"
)

type HttpCode int

type jsonOption struct {
	HttpCode int
	JsonType int
}
type optionFunc func(*jsonOption)

// http状态码
func WithHttpCode(httpCode int) optionFunc {
	return func(o *jsonOption) {
		if httpCode == 0 {
			httpCode = 200
		}
		o.HttpCode = httpCode
	}
}

// json序列化类型
func WithJsonType(jsonType int) optionFunc {
	return func(o *jsonOption) {
		if jsonType == 0 {
			jsonType = 1
		}
		o.JsonType = jsonType
	}
}

// JSON http json格式response响应
func JSON(c *gin.Context, v interface{}, e error, opts ...optionFunc) {
	v = convertListToEmptyArray(v)
	// json序列化选项
	jo := &jsonOption{
		HttpCode: http.StatusOK,
		JsonType: 1,
	}
	for _, opt := range opts {
		opt(jo)
	}
	// 序列化函数
	var jsonFunc func(code int, obj any)
	switch jo.JsonType {
	case 1:
		jsonFunc = c.JSON
	case 2:
		jsonFunc = c.PureJSON
	}
	// 状态码翻译
	acceptLanguage := strings.Split(c.GetHeader("Accept-Language"), ",")[0]
	if acceptLanguage == "" {
		acceptLanguage = "en-US"
	}
	// 正常的返回
	if e == nil {
		msg := i18n.T(acceptLanguage, 0)
		jsonFunc(jo.HttpCode, gin.H{
			"code":      0,
			"errorCode": 0,
			"msg":       msg,
			"data":      v,
		})
	}
	// 错误返回
	if e != nil {
		errorx, ok := e.(errorx.Error)
		// 自定义错误返回
		if ok {
			msg := i18n.T(acceptLanguage, errorx.Code)
			if msg == "" {
				msg = e.Error()
			}
			jsonFunc(jo.HttpCode, gin.H{
				"code":      errorx.Code,
				"errorCode": errorx.Code,
				"msg":       msg,
				"data":      v,
			})
		} else {
			// 系统错误
			jsonFunc(jo.HttpCode, gin.H{
				"code":      -1,
				"errorCode": -1,
				"msg":       e.Error(),
				"data":      v,
			})
		}
	}
}

// convertListToEmptyArray 当返回的list为空时，将list转为空数组
func convertListToEmptyArray(data interface{}) interface{} {
	if v, ok := data.(gin.H); ok {
		if v1, ok := v["list"]; ok {
			if reflect.TypeOf(v1).Kind() == reflect.Slice {
				if reflect.ValueOf(v1).Len() == 0 {
					v["list"] = []string{}
					data = v
				}
			}
		}
	} else if v, ok := data.(map[string]interface{}); ok {
		if v1, ok := v["list"]; ok {
			if reflect.TypeOf(v1).Kind() == reflect.Slice {
				if reflect.ValueOf(v1).Len() == 0 {
					v["list"] = []string{}
					data = v
				}
			}
		}
	}
	return data
}
