package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	ignorePaths map[string]struct{}
}

func LoginMiddleware(ignorePaths []string) *LoginMiddlewareBuilder {
	ignorePathsMap := make(map[string]struct{}, len(ignorePaths))
	for _, path := range ignorePaths {
		ignorePathsMap[path] = struct{}{}
	}
	return &LoginMiddlewareBuilder{
		ignorePaths: ignorePathsMap,
	}
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if _, ok := m.ignorePaths[path]; ok {
			// Do not check
			return
		}
		sess := sessions.Default(ctx)
		if sess.Get("userID") == nil {
			// abort the request and return error
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
