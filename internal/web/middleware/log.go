package middleware

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type LogMiddlewareBuilder struct {
	logFn         func(ctx context.Context, al AccessLog)
	allowReqBody  bool
	allowRespBody bool
	maxPathLen    int
	maxBodyLen    int
}

func NewLogMiddlewareBuilder(
	logFn func(ctx context.Context, al AccessLog),
) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFn:      logFn,
		maxPathLen: 1024,
		maxBodyLen: 2048,
	}
}

// }}}
// {{{ Other structs

type AccessLog struct {
	Path       string        `json:"path"`
	Method     string        `json:"method"`
	ReqBody    string        `json:"req_body"`
	RespBody   string        `json:"resp_body"`
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`
}

type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

// }}}
// {{{ Struct Methods

func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder {
	l.allowReqBody = true
	return l
}

func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder {
	l.allowRespBody = true
	return l
}

func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if len(path) > l.maxPathLen {
			path = path[:l.maxPathLen]
		}
		method := ctx.Request.Method
		al := AccessLog{
			Path:   path,
			Method: method,
		}
		if l.allowReqBody {
			// NOTE: Read from stream
			body, _ := ctx.GetRawData()
			if len(body) > l.maxBodyLen {
				al.ReqBody = string(body[:l.maxBodyLen])
			} else {
				al.ReqBody = string(body)
			}
			// NOTE: And put it back
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		start := time.Now()

		if l.allowRespBody {
			ctx.Writer = &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}

		defer func() {
			al.Duration = time.Since(start)
			l.logFn(ctx, al)
		}()

		ctx.Next()
	}
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.al.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
