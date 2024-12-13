package web

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTKey     = []byte("xQUPmbb2TP9CUyFZkgOnV3JQdr22ZNBx")
	RefreshKey = []byte("xQUPmbb2TP9CUyFZkgOnV3JQdr2fsNBx")
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserAgent string
	UID       int64
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UID int64
}

type jwtHandler struct {
	signingMethod jwt.SigningMethod
	refreshKey    []byte
}

func newJWTHandler() jwtHandler {
	return jwtHandler{
		signingMethod: jwt.SigningMethodHS256,
		refreshKey:    RefreshKey,
	}
}

func (h *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	err := h.setRefreshToken(ctx, uid)
	if err != nil {
		return err
	}
	tokenStr, err := h.generateJWTToken(ctx, uid)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h *jwtHandler) generateJWTToken(
	ctx *gin.Context,
	uid int64,
) (string, error) {
	uc := UserClaims{
		UID:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Second)),
		},
	}
	return h.makeToken(uc, JWTKey)
}

func (h *jwtHandler) setRefreshToken(ctx *gin.Context, uid int64) error {
	tokenStr, err := h.generateRefreshToken(uid)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (h *jwtHandler) generateRefreshToken(
	uid int64,
) (string, error) {
	rc := RefreshClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	return h.makeToken(rc, h.refreshKey)
}

func (h *jwtHandler) makeToken(claims jwt.Claims, signingKey []byte) (string, error) {
	token := jwt.NewWithClaims(h.signingMethod, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		slog.Error("token string generate error", "err", err)
		return "", err
	}
	return tokenStr, nil
}

func ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}

	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}

	return segs[1]
}
