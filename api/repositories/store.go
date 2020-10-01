package repositories

import (
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

type IDayRepository interface {
	Save(task model.Day) StoreResult
	GetDaysAsync(userId string, limit int, offset int) StoreChannel
	GetDayByIdAsync(id string) StoreChannel
	UpdateAsync(task model.Day) StoreChannel
	Delete(id string) StoreResult
}

type ITaskRepository interface {
	GetTaskAsync(model.TaskQuery) StoreChannel
	GetChildrenByTaskIdAsync(id string) StoreChannel
	GetTasksByParentAsync(id string) StoreChannel
	GetTaskByIdAsync(id string) StoreChannel
	Save(task model.Task) StoreResult
	UpdateAsync(task model.Task) StoreChannel
	Delete(id string) StoreResult
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
