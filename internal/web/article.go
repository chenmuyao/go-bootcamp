package web

import (
	"fmt"
	"time"

	"github.com/chenmuyao/generique/gslice"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type ArticleHandler struct {
	l   logger.Logger
	svc service.ArticleService
}

func NewArticleHandler(l logger.Logger, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		l:   l,
		svc: svc,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")

	g.POST("edit", ginx.WrapBodyAndClaims(h.l, h.Edit))
	g.POST("publish", ginx.WrapBodyAndClaims(h.l, h.Publish))
	g.POST("withdraw", ginx.WrapBodyAndClaims(h.l, h.Withdraw))

	// author
	g.GET("/detail/:id", ginx.WrapClaims(h.l, h.Detail))
	g.POST("/list", ginx.WrapBodyAndClaims(h.l, h.List))
}

func (h *ArticleHandler) Edit(
	ctx *gin.Context,
	req ArticleEditReq,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	aid, err := h.svc.Save(ctx, domain.Article{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			ID: uc.UID,
		},
	})
	switch err {
	case nil:
		return ginx.Result{
			Data: aid,
			Code: ginx.CodeOK,
		}, nil
	case service.ErrArticleNotFound:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "article not found",
		}, nil
	default:
		return ginx.InternalServerErrorResult, fmt.Errorf("Save article failed: %w", err)
	}
}

func (h *ArticleHandler) Publish(
	ctx *gin.Context,
	req ArticlePublishReq,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	aid, err := h.svc.Publish(ctx, domain.Article{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			ID: uc.UID,
		},
	})
	switch err {
	case nil:
		return ginx.Result{
			Data: aid,
			Code: ginx.CodeOK,
		}, nil
	case service.ErrArticleNotFound:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "article not found",
		}, nil
	default:
		return ginx.InternalServerErrorResult, fmt.Errorf("Publish article failed: %w", err)
	}
}

func (h *ArticleHandler) Withdraw(
	ctx *gin.Context,
	req ArticleWithdrawReq,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	err := h.svc.Withdraw(ctx, uc.UID, req.ID)
	switch err {
	case nil:
		return ginx.Result{
			Code: ginx.CodeOK,
		}, nil
	case service.ErrArticleNotFound:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "article not found",
		}, nil
	default:
		return ginx.InternalServerErrorResult, fmt.Errorf(
			"Withdraw article %d failed: %w",
			req.ID,
			err,
		)
	}
}

func (h *ArticleHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	return ginx.Result{}, nil
}

func (h *ArticleHandler) List(
	ctx *gin.Context,
	page Page,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	articles, err := h.svc.GetByAuthor(ctx, uc.UID, page.Offset, page.Limit)
	switch err {
	case nil:
		return ginx.Result{
			Code: ginx.CodeOK,
			Data: gslice.Map(articles, func(id int, src domain.Article) ArticleVO {
				return ArticleVO{
					ID:      src.ID,
					Title:   src.Title,
					Content: src.Content,
					Status:  uint8(src.Status),
					Ctime:   src.Ctime.Format(time.DateTime),
					Utime:   src.Ctime.Format(time.DateTime),
				}
			}),
		}, nil
	case service.ErrArticleNotFound:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "article not found",
		}, nil
	default:
		return ginx.InternalServerErrorResult, &logger.LogError{
			Msg: "Get articles by author failed",
			Fields: []logger.Field{
				logger.Int64("uid", uc.UID),
				logger.Int("offset", page.Offset),
				logger.Int("limit", page.Limit),
			},
		}
	}
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
