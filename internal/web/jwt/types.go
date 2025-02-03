package jwt

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=./types.go -package=jwtmocks -destination=./mocks/jwt.mock.go
type Handler interface {
	ExtractToken(ctx *gin.Context) string
	ClearToken(ctx *gin.Context) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	GenerateRefreshToken(uid int64, ssid string) (string, error)
	GenerateJWTToken(ctx *gin.Context, uid int64, ssid string) (string, error)
	CheckSession(ctx *gin.Context, ssid string) error
}
