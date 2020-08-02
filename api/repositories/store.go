package repositories

import (
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/Banyango/gifoody_server/api/repositories/stores"
	. "github.com/Banyango/gifoody_server/api/repositories/util"
	"github.com/jmoiron/sqlx"
)

type DBContext interface {
	Post() IPostRepository
}

type IPostRepository interface {
	FindPosts(model.PostQuery) StoreChannel
}

type IUserRepository interface {
	GetUserByEmailAsync(email string) StoreChannel
	GetUserByConfirmTokenAsync(token string) StoreChannel
	GetUserByIdAsync(id string) StoreChannel
	Save(user model.User) StoreResult
	UpdateAsync(user model.User) StoreChannel
	DeleteForgotUser(id string) StoreResult
	SaveForgotUser(user model.ForgotUser) StoreResult
}

type AppStore struct {
	db   *sqlx.DB
	post IPostRepository
	user IUserRepository
}

func (self *AppStore) Post() IPostRepository {
	return self.post
}

func (self *AppStore) User() IUserRepository {
	return self.user
}

func NewAppStore(db *sqlx.DB) *AppStore {
	store := new(AppStore)

	store.post = stores.NewPostSQLStore(db)
	store.user = stores.NewUserSQLStore(db)

	return store
}
