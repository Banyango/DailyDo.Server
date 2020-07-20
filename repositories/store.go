package repositories

import (
	"github.com/Banyango/gifoody_server/model"
	. "github.com/Banyango/gifoody_server/repositories/sql"
	. "github.com/Banyango/gifoody_server/repositories/util"
	"github.com/jmoiron/sqlx"
)

type DBContext interface {
	Post() IPostRepository
}

type IPostRepository interface {
	FindPosts(model.PostQuery) StoreChannel
}

type IUserRepository interface {
	GetUserByEmail(email string) StoreChannel
	GetUserByConfirmToken(token string) StoreChannel
	GetUserById(id string) StoreChannel
	Save(user model.User) StoreChannel
	Update(user model.User) StoreChannel
	DeleteForgotUser(id string) StoreResult
	SaveForgotUser(user model.ForgotUser) StoreChannel
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

	store.post = NewPostSQLStore(db)
	store.user = NewUserSQLStore(db)

	return store
}
