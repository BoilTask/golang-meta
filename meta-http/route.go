package metahttp

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	metacontroller "meta/controller"
	"reflect"
	"strings"
)

var AuthMiddleware gin.HandlerFunc
var AuthMiddlewareOptional gin.HandlerFunc

func AutoRegisterRoute(r gin.IRoutes, basePath string, controller metacontroller.Interface, authMiddlewareType AuthMiddlewareType) {
	switch authMiddlewareType {
	case AuthMiddlewareTypeNone:
		AutoRegisterRouteWithHandle(r, basePath, controller)
	case AuthMiddlewareTypeOptional:
		if AuthMiddlewareOptional == nil {
			AutoRegisterRouteWithHandle(r, basePath, controller, AuthMiddleware)
		} else {
			AutoRegisterRouteWithHandle(r, basePath, controller, AuthMiddlewareOptional)
		}
	default:
		AutoRegisterRouteWithHandle(r, basePath, controller, AuthMiddleware)
	}
}

func AutoRegisterRouteWithHandle(r gin.IRoutes, basePath string, controller metacontroller.Interface, handlers ...gin.HandlerFunc) {
	ctrlVal := reflect.ValueOf(controller)
	ctrlType := reflect.TypeOf(controller)

	for i := 0; i < ctrlVal.NumMethod(); i++ {
		method := ctrlVal.Method(i)
		methodName := ctrlType.Method(i).Name

		// 根据方法名自动推断 HTTP方法和路径
		var maps = map[string]string{
			"Get":     "GET",
			"Post":    "POST",
			"Put":     "PUT",
			"Delete":  "DELETE",
			"Options": "OPTIONS",
		}
		for k, v := range maps {
			if strings.HasPrefix(methodName, k) {
				httpMethod := v
				path := strings.TrimPrefix(methodName, k)

				path = strings.Join(splitCamelCase(path), "/")

				realPath := basePath
				if path != "" {
					realPath = realPath + "/" + path
				}

				callHandler := func(c *gin.Context) {
					method.Call([]reflect.Value{reflect.ValueOf(c)})
				}

				if len(handlers) > 0 {
					realHandlers := append(handlers, callHandler)
					r.Handle(httpMethod, realPath, realHandlers...)
				} else {
					r.Handle(httpMethod, realPath, callHandler)
				}

				slog.Info(
					"auto register path",
					"method",
					httpMethod,
					"path",
					realPath,
					"extraHandler",
					len(handlers),
				)
			}
		}
	}
}

func splitCamelCase(s string) []string {
	var words []string
	var currentWord strings.Builder
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			words = append(words, strings.ToLower(currentWord.String()))
			currentWord.Reset()
		}
		currentWord.WriteRune(char)
	}
	words = append(words, strings.ToLower(currentWord.String()))
	return words
}
