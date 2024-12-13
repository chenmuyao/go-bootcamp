package web

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWTKey = []byte("xQUPmbb2TP9CUyFZkgOnV3JQdr22ZNBx")

type UserClaims struct {
	jwt.RegisteredClaims
	UserAgent string
	UID       int64
}

type jwtHandler struct{}

func (h *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	tokenStr, err := h.generateJWTToken(ctx, uid)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h *jwtHandler) generateJWTToken(ctx *gin.Context, uid int64) (string, error) {
	uc := UserClaims{
		UID:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		slog.Error("token string generate error", "err", err)
		return "", err
	}
	return tokenStr, nil
}
