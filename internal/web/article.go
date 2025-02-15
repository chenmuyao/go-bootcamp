package web

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/chenmuyao/generique/gslice"
	intrDomain "github.com/chenmuyao/go-bootcamp/interactive/domain"
	intrService "github.com/chenmuyao/go-bootcamp/interactive/service"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type ArticleHandler struct {
	l       logger.Logger
	svc     service.ArticleService
	intrSvc intrService.InteractiveService
	biz     string
}

func NewArticleHandler(
	l logger.Logger,
	svc service.ArticleService,
	intrSvc intrService.InteractiveService,
) *ArticleHandler {
	return &ArticleHandler{
		l:       l,
		svc:     svc,
		intrSvc: intrSvc,
		biz:     "article",
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
	// normally: /list?offset=?&limit=?
	g.POST("/list", ginx.WrapBodyAndClaims(h.l, h.List))

	// get published article (reader)
	pub := g.Group("/pub")
	pub.GET("/:id", ginx.WrapClaims(h.l, h.PubDetail))
	pub.GET("/top_like", ginx.WrapLog(h.l, h.TopLike))
	// True: like; False: cancel like
	pub.POST("/like", ginx.WrapBodyAndClaims(h.l, h.Like))
	pub.POST("/collect", ginx.WrapBodyAndClaims(h.l, h.Collect))
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
		return ginx.InternalServerErrorResult, logger.LError(
			"Save article failed: %w",
			logger.Error(err),
		)
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.l.Warn("wrong id", logger.String("id", idStr), logger.Error(err))
		return ginx.InternalServerErrorResult, nil
	}
	article, err := h.svc.GetByID(ctx, id)
	if err != nil {
		return ginx.InternalServerErrorResult, fmt.Errorf(
			"Get article %d detail failed: %w",
			id,
			err,
		)
	}
	if article.Author.ID != uc.UID {
		return ginx.Result{
				Code: ginx.CodeUserSide,
				Msg:  "article not found",
			}, logger.LError("invalid article query",
				logger.Int64("id", id),
				logger.Int64("uid", uc.UID),
			)
	}
	return ginx.Result{
		Code: ginx.CodeOK,
		Data: ArticleVO{
			ID:    article.ID,
			Title: article.Title,
			// Abstract: article.Abstract(),
			Content: article.Content,
			Status:  uint8(article.Status),
			Ctime:   article.Ctime.Format(time.DateTime),
			Utime:   article.Ctime.Format(time.DateTime),
		},
	}, nil
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
					ID:       src.ID,
					Title:    src.Title,
					Abstract: src.Abstract(),
					// Content:  src.Content,
					Status: uint8(src.Status),
					Ctime:  src.Ctime.Format(time.DateTime),
					Utime:  src.Ctime.Format(time.DateTime),
				}
			}),
		}, nil
	case service.ErrArticleNotFound:
		return ginx.Result{
			Code: ginx.CodeUserSide,
			Msg:  "article not found",
		}, nil
	default:
		return ginx.InternalServerErrorResult,
			logger.LError("Get articles by author failed",
				logger.Int64("uid", uc.UID),
				logger.Int("offset", page.Offset),
				logger.Int("limit", page.Limit),
			)
	}
}

func (h *ArticleHandler) PubDetail(
	ctx *gin.Context,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.l.Warn("wrong id", logger.String("id", idStr), logger.Error(err))
		return ginx.InternalServerErrorResult, nil
	}

	var (
		eg      errgroup.Group
		article domain.Article
		intr    intrDomain.Interactive
	)

	eg.Go(func() error {
		var er error
		article, er = h.svc.GetPubByID(ctx, id, uc.UID)
		return er
	})

	eg.Go(func() error {
		var er error
		intr, er = h.intrSvc.Get(ctx, h.biz, id, uc.UID)
		h.l.Debug("intr", logger.Field{Key: "intr", Value: intr}, logger.Error(er))

		return er
	})

	err = eg.Wait()
	if err != nil {
		return ginx.InternalServerErrorResult, logger.LError(
			"Get published article detail failed",
			logger.Int64("id", id),
			logger.Int64("uid", uc.UID),
			logger.Error(err),
		)
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := h.intrSvc.IncrReadCnt(ctx, h.biz, article.ID)
		if er != nil {
			h.l.Error(
				"failed to update read count",
				logger.String("biz", h.biz),
				logger.Int64("bizID", article.ID),
				logger.Error(er),
			)
		}
	}()

	return ginx.Result{
		Code: ginx.CodeOK,
		Data: ArticleVO{
			ID:    article.ID,
			Title: article.Title,
			// Abstract: article.Abstract(),
			Content:    article.Content,
			AuthorID:   article.Author.ID,
			AuthorName: article.Author.Name,
			Status:     uint8(article.Status),
			Ctime:      article.Ctime.Format(time.DateTime),
			Utime:      article.Ctime.Format(time.DateTime),

			ReadCnt:    intr.ReadCnt,
			LikeCnt:    intr.LikeCnt,
			CollectCnt: intr.CollectCnt,
			Liked:      intr.Liked,
			Collected:  intr.Collected,
		},
	}, nil
}

func (h *ArticleHandler) TopLike(ctx *gin.Context) (ginx.Result, error) {
	// Get top limit
	var limit int
	limitStr := ctx.Query("limit")
	if res, err := strconv.Atoi(limitStr); err != nil {
		limit = res
	}

	articleIDs, err := h.intrSvc.GetTopLike(ctx, h.biz, limit)
	if err != nil {
		return ginx.InternalServerErrorResult, logger.LError(
			"failed to get top like",
			logger.Error(err),
		)
	}

	articles, err := h.svc.BatchGetPubByIDs(ctx, articleIDs)
	if err != nil {
		return ginx.InternalServerErrorResult, logger.LError(
			"failed to get articles",
			logger.String("biz", h.biz),
			logger.Error(err),
		)
	}

	intrs, err := h.intrSvc.MustBatchGet(ctx, h.biz, articleIDs)
	if err != nil {
		return ginx.InternalServerErrorResult, logger.LError(
			"failed to get interactives",
			logger.String("biz", h.biz),
			logger.Error(err),
		)
	}

	intrArticles := gslice.Map(articles, func(id int, src domain.Article) ArticleVO {
		return ArticleVO{
			ID:         src.ID,
			Title:      src.Title,
			Abstract:   src.Abstract(),
			AuthorID:   src.Author.ID,
			AuthorName: src.Author.Name,
			Ctime:      src.Ctime.Format(time.DateTime),
			Utime:      src.Ctime.Format(time.DateTime),
			LikeCnt:    intrs[id].LikeCnt,
		}
	})
	return ginx.Result{
		Code: ginx.CodeOK,
		Data: intrArticles,
	}, nil
}

func (h *ArticleHandler) Like(ctx *gin.Context, req Like, uc ijwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		err = h.intrSvc.Like(ctx, h.biz, req.ID, uc.UID)
	} else {
		err = h.intrSvc.CancelLike(ctx, h.biz, req.ID, uc.UID)
	}
	if err != nil {
		return ginx.InternalServerErrorResult, logger.LError(
			"failed to like or cancel like",
			logger.Error(err),
			logger.Int64("uid", uc.UID),
			logger.Int64("aid", req.ID),
		)
	}
	return ginx.Result{
		Code: ginx.CodeOK,
	}, nil
}

func (h *ArticleHandler) Collect(
	ctx *gin.Context,
	req Collect,
	uc ijwt.UserClaims,
) (ginx.Result, error) {
	var err error
	if req.Collected {
		err = h.intrSvc.Collect(ctx, h.biz, req.ID, req.CID, uc.UID)
	} else {
		err = h.intrSvc.CancelCollect(ctx, h.biz, req.ID, req.CID, uc.UID)
	}
	if err != nil {
		return ginx.InternalServerErrorResult, logger.LError(
			"failed to collect article",
			logger.Error(err),
			logger.Int64("uid", uc.UID),
			logger.Int64("aid", req.ID),
		)
	}
	return ginx.Result{
		Code: ginx.CodeOK,
	}, nil
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
