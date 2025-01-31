package middleware

import (
	"net/http"

	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type LoginJWT struct {
	ignorePaths map[string]struct{}
	ijwt.Handler
}

func NewLoginJWTBuilder(hdl ijwt.Handler, ignorePaths []string) *LoginJWT {
	ignorePathsMap := make(map[string]struct{}, len(ignorePaths))
	for _, path := range ignorePaths {
		ignorePathsMap[path] = struct{}{}
	}
	return &LoginJWT{
		ignorePaths: ignorePathsMap,
		Handler:     hdl,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (m *LoginJWT) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if _, ok := m.ignorePaths[path]; ok {
			// Do not check
			return
		}
		// Authorization: Bearer XXXX
		tokenStr := m.ExtractToken(ctx)

		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(t *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})
		if err != nil {
			// token cannot be parsed
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			// parsed but unauthorized
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// if uc.UserAgent != ctx.GetHeader(httpx.UserAgent) {
		// 	// NOTE: Instrument here. Might be attackers.
		// 	// A better option is to use the browser's fingerprint.
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		// Check if logged out
		if err := m.CheckSession(ctx, uc.SSID); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// NOTE: Automatic refresh
		// expireTime := uc.ExpiresAt.Time
		// if time.Until(expireTime) < 29*time.Minute {
		// 	// refresh every minute
		// 	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * time.Minute))
		// 	tokenStr, err = token.SignedString(web.JWTKey)
		// 	ctx.Header("x-jwt-token", tokenStr)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// }

		ctx.Set("user", uc)
	}
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
