package repositories

import (
	"context"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/Banyango/gifoody_server/api/repositories/stores"
	. "github.com/Banyango/gifoody_server/api/repositories/util"
	"github.com/jmoiron/sqlx"
)

type DBContext interface {
	Task() ITaskRepository
	Day() IDayRepository
	User() IUserRepository
}

type ISQLStore interface {
	Execute(ctx context.Context, fn func(c context.Context) error) error
}

type IDayRepository interface {
	Save(task model.Day) StoreResult
	GetDaysAsync(userId string, limit int, offset int) StoreChannel
	GetDayByIdAsync(id string) StoreChannel
	UpdateAsync(task model.Day) StoreChannel
	Delete(id string) StoreResult
}

type ITaskRepository interface {
	ISQLStore
	GetTaskAsync(query model.TaskQuery, ctx context.Context) StoreChannel
	GetChildrenByTaskIdAsync(id string, ctx context.Context) StoreChannel
	GetTasksByParentAsync(id string, ctx context.Context) StoreChannel
	GetTaskByIdAsync(id string, ctx context.Context) StoreChannel
	GetMaxOrder(id string, ctx context.Context) StoreChannel
	GetTaskByOrderIdAsync(parent string, id string, ctx context.Context) StoreChannel
	Save(task model.Task, ctx context.Context) StoreResult
	UpdateAsync(task model.Task, ctx context.Context) StoreChannel
	Delete(id string, ctx context.Context) StoreResult
}

type IUserRepository interface {
	GetUserByEmailOrUsernameAsync(email string, username string) StoreChannel
	GetUserByConfirmTokenAsync(token string) StoreChannel
	GetForgotUserByTokenAsync(token string) StoreChannel
	GetUserByIdAsync(id string) StoreChannel
	Save(user model.User) StoreResult
	UpdateAsync(user model.User) StoreChannel
	DeleteForgotUser(id string) StoreResult
	SaveForgotUser(user model.ForgotUser) StoreResult
}

type AppStore struct {
	db   *sqlx.DB
	task ITaskRepository
	user IUserRepository
	day  IDayRepository
}

func (self *AppStore) Task() ITaskRepository {
	return self.task
}

func (self *AppStore) User() IUserRepository {
	return self.user
}

func (self *AppStore) Day() IDayRepository {
	return self.day
}

func NewAppStore(db *sqlx.DB) *AppStore {
	store := new(AppStore)

	store.day = stores.NewDaySQLStore(db)
	store.task = stores.NewTaskSQLStore(db)
	store.user = stores.NewUserSQLStore(db)

	return store
}
