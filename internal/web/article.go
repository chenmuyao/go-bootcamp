package web

import (
	"fmt"

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

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
