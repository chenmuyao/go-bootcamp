package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

var (
	ErrCodeSendTooMany   = errors.New("too frequent send requests")
	ErrCodeVerifyTooMany = errors.New("too frequent verify requests")
	ErrKeyNotExist       = errors.New("inexisted key")
)

// }}}
// {{{ Interface

//go:generate mockgen -source=./types.go -package=cachemocks -destination=./mocks/cache.mock.go

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
	// Key(biz, phone string) string
}

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
	BatchGet(ctx context.Context, uids []int64) ([]domain.User, error)
	BatchSet(ctx context.Context, users []domain.User) error
}

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, article domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, article domain.Article) error
	BatchGetPub(ctx context.Context, ids []int64) ([]domain.Article, error)
	BatchSetPub(ctx context.Context, articles []domain.Article) error
}

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

// }}}
// {{{ Struct

type BaseCodeCache struct{}

type BaseUserCache struct{}

type BaseArticleCache struct{}

type BaseInteractiveCache struct{}

type BaseRankingCache struct{}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (b *BaseCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *BaseUserCache) Key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (c *BaseArticleCache) Key(motif string, uid int64) string {
	return fmt.Sprintf("article:%s:%d", motif, uid)
}

func (c *BaseInteractiveCache) Key(biz string, bizID int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizID)
}

func (c *BaseRankingCache) Key(biz string) string {
	return fmt.Sprintf("ranking:%s", biz)
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
