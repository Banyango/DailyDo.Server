package repositories

import (
	"database/sql"
	"github.com/Banyango/gifoody_server/model"
	. "github.com/Banyango/gifoody_server/repositories/sql"
	. "github.com/Banyango/gifoody_server/repositories/util"
)

type DBContext interface {
	Post() IPostRepository
}

type IPostRepository interface {
	FindPosts(model.PostQuery) StoreChannel
}

type AppStore struct {
	db   *sql.DB
	post IPostRepository
}

func (self *AppStore) Post() IPostRepository {
	return self.post;
}

func NewAppStore(db *sql.DB) *AppStore {
	store := new(AppStore)

	store.post = NewPostSQLStore(db)

	return store
}
