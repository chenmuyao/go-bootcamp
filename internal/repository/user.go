package repository

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"time"

	"github.com/chenmuyao/generique/gslice"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

var (
	ErrDuplicatedUser = dao.ErrDuplicatedUser
	// NOTE: Strongly related to the service
	ErrUserNotFound = dao.ErrRecordNotFound
)

// }}}
// {{{ Interface

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, userID int64) (domain.User, error)
	BatchFindByIDs(ctx context.Context, userIDs []int64) ([]domain.User, error)
}

// }}}
// {{{ Struct

type CachedUserRepository struct {
	cache cache.UserCache
	dao   dao.UserDAO
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	daoUser, err := repo.dao.Insert(ctx, repo.userDomainToDAO(&u))
	if err != nil {
		return domain.User{}, err
	}
	return repo.userDAOToDomain(&daoUser), nil
}

func (repo *CachedUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.userDAOToDomain(&u), nil
}

func (repo *CachedUserRepository) FindByPhone(
	ctx context.Context,
	phone string,
) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.userDAOToDomain(&u), nil
}

func (repo *CachedUserRepository) UpdateProfile(ctx context.Context, user *domain.User) error {
	err := repo.dao.UpdateProfile(ctx, repo.userDomainToDAO(user))
	if err != nil {
		return err
	}
	return nil
}

// BatchFindByIDs implements UserRepository.
func (repo *CachedUserRepository) BatchFindByIDs(
	ctx context.Context,
	userIDs []int64,
) ([]domain.User, error) {
	du, err := repo.cache.BatchGet(ctx, userIDs)
	// found, return
	if err == nil {
		return du, nil
	}
	slog.Error("redis get", "err", err)

	daoUsers, err := repo.dao.BatchFindByIDs(ctx, userIDs)
	if err != nil {
		return []domain.User{}, err
	}

	domainUsers := gslice.Map(daoUsers, func(id int, src dao.User) domain.User {
		return repo.userDAOToDomain(&src)
	})

	go func() {
		err = repo.cache.BatchSet(ctx, domainUsers)
		if err != nil {
			// Network, or redis crash
			slog.Error("redis set", "err", err)
		}
	}()

	return du, nil
}

// NOTE: Ok for normal case. But if cache penetration happens, the DB can be crashed
// by queries
func (repo *CachedUserRepository) FindByID(ctx context.Context, userID int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, userID)
	// found, return
	if err == nil {
		return du, nil
	}
	slog.Error("redis get", "err", err)

	// err != nil
	// 1. key inexistant: Redis ok
	// 2. Redis nok (network or redis crash) ==> cache penetration
	u, err := repo.dao.FindByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	du = repo.userDAOToDomain(&u)

	// NOTE: Async 1ms economy
	go func() {
		err = repo.cache.Set(ctx, du)
		if err != nil {
			// Network, or redis crash
			slog.Error("redis set", "err", err)
		}
	}()

	return du, nil
}

// NOTE: Conservative method, for tasks with low priority. This helps keep the
// database working for other services.
func (repo *CachedUserRepository) FindByIdV1(
	ctx context.Context,
	userID int64,
) (domain.User, error) {
	du, err := repo.cache.Get(ctx, userID)
	// found, return
	switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExist:
		// 1. key inexistant: Redis ok
		u, err := repo.dao.FindByID(ctx, userID)
		if err != nil {
			return domain.User{}, err
		}

		du = repo.userDAOToDomain(&u)

		// NOTE: Async 1ms economy
		go func() {
			err = repo.cache.Set(ctx, du)
			if err != nil {
				// Network, or redis crash
				log.Println(err)
			}
		}()
		return du, nil
	default:
		// 2. Redis nok (network or redis crash) ==>
		// Degrade: avoid cache penetration
		return domain.User{}, err
	}
}

func (repo *CachedUserRepository) userDAOToDomain(u *dao.User) domain.User {
	return domain.User{
		ID:       u.ID,
		Password: u.Password,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Name:     u.Name,
		Birthday: time.UnixMilli(u.Birthday),
		Profile:  u.Profile,
	}
}

func (repo *CachedUserRepository) userDomainToDAO(u *domain.User) dao.User {
	return dao.User{
		ID:       u.ID,
		Password: u.Password,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Phone:    sql.NullString{String: u.Phone, Valid: u.Phone != ""},
		Name:     u.Name,
		Birthday: u.Birthday.UnixMilli(),
		Profile:  u.Profile,
	}
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
