package metahttp

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log/slog"
	metaerrorcode "meta/error-code"
	"meta/meta-flag"
	"meta/meta-log"
	metapanic "meta/meta-panic"
	metaresponse "meta/meta-response"
	"meta/routine"
	"meta/subsystem"
	"net/http"
	"time"
)

type Subsystem struct {
	subsystem.Subsystem
	GetPort    func() int32
	ProcessGin func(r *gin.Engine)
}

func (s *Subsystem) GetName() string {
	return "Http"
}

func (s *Subsystem) Start() error {
	routine.SafeGoWithRestart(
		"Http start",
		func() error {
			err := s.startSubsystem()
			if err != nil {
				return err
			}
			return nil
		},
	)
	return nil
}

func (s *Subsystem) startSubsystem() error {

	port := s.GetPort()

	if !metaflag.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建默认路由引擎
	r := gin.New()

	// 手动添加中间件
	r.Use(GinLogger(metalog.GetLogger()))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// 需要先设置CORS，否则错误可能会被拦截
	r.Use(CORSMiddleware())
	r.Use(ErrorHandlingMiddleware())
	r.NoRoute(
		func(ctx *gin.Context) {
			metaresponse.NewResponse(ctx, http.StatusNotFound)
		},
	)

	s.ProcessGin(r)

	slog.Info("Http server listen start", "port", port)

	// 启动
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	return nil
}

func GinLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 处理请求
		ctx.Next()

		// 记录结束时间和持续时间
		duration := time.Since(start)

		// 记录日志信息
		logger.Info(
			"Gin",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.String("client_ip", ctx.ClientIP()),
			slog.Int("status_code", ctx.Writer.Status()),
			slog.String("user_agent", ctx.Request.UserAgent()),
			slog.Duration("latency", duration),
		)
	}
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response := metaresponse.Response{
					Code: int(metaerrorcode.PanicError),
				}
				ctx.JSON(http.StatusOK, response)
				metapanic.ProcessPanic(
					"gin panic", err,
					"ErrorHandlingMiddleware panic\nip: %s\nhost: %s\npath: %s\nmethod: %s\nparam: %v",
					ctx.ClientIP(),
					ctx.Request.Host,
					ctx.Request.URL.Path,
					ctx.Request.Method,
					ctx.Request.URL.Query(),
				)
				ctx.Abort()
			}
			for _, err := range ctx.Errors {
				metapanic.ProcessError(
					err.Err,
					"ErrorHandlingMiddleware error\nip: %s\nhost: %s\npath: %s\nmethod: %s\nparam: %v",
					ctx.ClientIP(),
					ctx.Request.Host,
					ctx.Request.URL.Path,
					ctx.Request.Method,
					ctx.Request.URL.Query().Encode(),
				)
			}

		}()
		ctx.Next()
	}
}
