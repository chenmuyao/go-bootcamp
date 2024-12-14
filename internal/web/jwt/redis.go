package jwt

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/consts"
	"github.com/chenmuyao/go-bootcamp/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// {{{ Consts
// }}}
// {{{ Global Varirables

var (
	JWTKey     = []byte("xQUPmbb2TP9CUyFZkgOnV3JQdr22ZNBx")
	RefreshKey = []byte("xQUPmbb2TP9CUyFZkgOnV3JQdr2fsNBx")
	ssidPrefix = "users:ssid"
)

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type RedisJWTHandler struct {
	signingMethod jwt.SigningMethod
	client        redis.Cmdable
	rcExpiration  time.Duration
	jwtExpiration time.Duration
}

var _ Handler = &RedisJWTHandler{}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		signingMethod: jwt.SigningMethodHS256,
		client:        client,
		rcExpiration:  time.Hour * 24 * 7,
		jwtExpiration: 30 * time.Minute,
	}
}

// }}}
// {{{ Other structs

type UserClaims struct {
	jwt.RegisteredClaims
	UserAgent string
	UID       int64
	SSID      string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UID  int64
	SSID string
}

// }}}
// {{{ Struct Methods

func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader(httpx.Authorization)
	if authCode == "" {
		return ""
	}

	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}

	return segs[1]
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header(consts.XJWTToken, "")
	ctx.Header(consts.XRefreshToken, "")
	uc := ctx.MustGet("user").(UserClaims)
	return h.client.Set(ctx, fmt.Sprintf("%s:%s", ssidPrefix, uc.SSID), "", h.rcExpiration).Err()
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisJWTHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	tokenStr, err := h.GenerateJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	ctx.Header(consts.XJWTToken, tokenStr)
	return nil
}

func (h *RedisJWTHandler) GenerateJWTToken(
	ctx *gin.Context,
	uid int64,
	ssid string,
) (string, error) {
	uc := UserClaims{
		UID:       uid,
		UserAgent: ctx.GetHeader(httpx.UserAgent),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.jwtExpiration)),
		},
		SSID: ssid,
	}
	return h.makeToken(uc, JWTKey)
}

func (h *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	tokenStr, err := h.GenerateRefreshToken(uid, ssid)
	if err != nil {
		return err
	}
	ctx.Header(consts.XRefreshToken, tokenStr)
	return nil
}

func (h *RedisJWTHandler) GenerateRefreshToken(uid int64, ssid string) (string, error) {
	rc := RefreshClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
		SSID: ssid,
	}
	return h.makeToken(rc, RefreshKey)
}

func (h *RedisJWTHandler) makeToken(claims jwt.Claims, signingKey []byte) (string, error) {
	token := jwt.NewWithClaims(h.signingMethod, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		slog.Error("token string generate error", "err", err)
		return "", err
	}
	return tokenStr, nil
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	// NOTE: Ignore redis error. In case of redis error, most users can still use the service
	cnt, err := h.client.Exists(ctx, fmt.Sprintf("%s:%s", ssidPrefix, ssid)).Result()
	if err != nil {
		// warning
		slog.Warn("Redis error", "err", err)
	}
	if cnt > 0 {
		// invalid token
		return errors.New("invalide token")
	}
	return nil
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
