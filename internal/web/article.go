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
}

func (h *ArticleHandler) Edit(
	ctx *gin.Context,
	req ArticleEditReq,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	aid, err := h.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			ID: uc.UID,
		},
	})
	if err != nil {
		return ginx.InternalServerErrorResult, fmt.Errorf("Save article failed: %w", err)
	}
	return ginx.Result{
		Data: aid,
		Code: ginx.CodeOK,
	}, nil
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
