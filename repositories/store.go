package repositories

import "github.com/Banyango/gifoody_server/model"

type StoreResult struct {
	Data interface{}
	Err  error
}

type StoreChannel chan StoreResult

type DBContext interface {
	Post() IPostRepository
}

type IPostRepository interface {
	FindPosts(model.PostQuery) StoreChannel
}
