package metahttp

type AuthMiddlewareType int

var (
	AuthMiddlewareTypeNone     AuthMiddlewareType = 0
	AuthMiddlewareTypeRequire  AuthMiddlewareType = 1
	AuthMiddlewareTypeOptional AuthMiddlewareType = 2
)
