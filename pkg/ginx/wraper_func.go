package ginx

import (
	"net/http"
	"strconv"

	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
)

var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{
		"code",
	})
	if err := prometheus.Register(vector); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			vector = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
}

func WrapBodyAndClaims[Req any, Claims jwt.Claims](
	l logger.Logger,
	bizFn func(ctx *gin.Context, req Req, uc Claims) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req

		if err := ctx.Bind(&req); err != nil {
			l.Error("user input error", logger.Error(err))
			return
		}

		l.Debug("request params", logger.Field{
			Key:   "req",
			Value: req,
		})

		val, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		uc, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		res, err := bizFn(ctx, req, uc)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if logError, ok := err.(*logger.LogError); ok {
			l.Error(logError.Msg, logError.Fields...)
		} else if err != nil {
			l.Error("failed to handle request", logger.Error(err))
		}

		// if result is empty, then we consider that it has already been set
		// into the response. just return
		var empty Result
		if res != empty {
			ctx.JSON(ResultToStatus(res), res)
		}
	}
}

func WrapBody[Req any](
	l logger.Logger,
	bizFn func(ctx *gin.Context, req Req) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req

		if err := ctx.Bind(&req); err != nil {
			l.Error("user input error", logger.Error(err))
			return
		}

		l.Debug("request params", logger.Field{
			Key:   "req",
			Value: req,
		})

		res, err := bizFn(ctx, req)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if logError, ok := err.(*logger.LogError); ok {
			l.Error(logError.Msg, logError.Fields...)
		} else if err != nil {
			l.Error("failed to handle request", logger.Error(err))
		}

		var empty Result
		if res != empty {
			ctx.JSON(ResultToStatus(res), res)
		}
	}
}

func WrapClaims[Claims jwt.Claims](
	l logger.Logger,
	bizFn func(ctx *gin.Context, uc Claims) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		uc, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		res, err := bizFn(ctx, uc)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if logError, ok := err.(*logger.LogError); ok {
			l.Error(logError.Msg, logError.Fields...)
		} else if err != nil {
			l.Error("failed to handle request", logger.Error(err))
		}

		var empty Result
		if res != empty {
			l.Debug("res", logger.Field{Key: "res", Value: res})
			ctx.JSON(ResultToStatus(res), res)
		}
	}
}

func WrapLog(l logger.Logger, bizFn func(ctx *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := bizFn(ctx)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if logError, ok := err.(*logger.LogError); ok {
			l.Error(logError.Msg, logError.Fields...)
		} else if err != nil {
			l.Error("failed to handle request", logger.Error(err))
		}

		var empty Result
		if res != empty {
			ctx.JSON(ResultToStatus(res), res)
		}
	}
}
