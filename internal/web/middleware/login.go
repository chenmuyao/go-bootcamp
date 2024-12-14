package middleware

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

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

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if _, ok := m.ignorePaths[path]; ok {
			// Do not check
			return
		}
		sess := sessions.Default(ctx)
		userID := sess.Get("userID")
		if userID == nil {
			// abort the request and return error
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Refresh session every minute
		now := time.Now()
		const updateTimeKey = "update_time"

		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Minute {
			sess.Set(updateTimeKey, now)
			sess.Options(sessions.Options{
				MaxAge: 900,
			})
			err := sess.Save()
			if err != nil {
				log.Println(err)
			}
			ctx.Next()
		}
	}
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
