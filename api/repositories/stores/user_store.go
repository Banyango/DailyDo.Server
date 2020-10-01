package stores

import (
	"fmt"
	. "github.com/Banyango/gifoody_server/api/model"
	. "github.com/Banyango/gifoody_server/api/repositories/util"
	"github.com/jmoiron/sqlx"
)

type UserSQLStore struct {
	db *sqlx.DB
}

func NewUserSQLStore(session *sqlx.DB) *UserSQLStore {
	self := new(UserSQLStore)

	self.db = session

	return self
}

func (self *UserSQLStore) GetUserByEmailOrUsernameAsync(email string, username string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := User{}
		err := self.db.Get(&user, "SELECT * from users WHERE email = ? OR username = ?", email, username)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) GetUserByConfirmTokenAsync(token string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := User{}
		err := self.db.Get(&user, "SELECT * from users WHERE confirm_token = ?", token)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) GetForgotUserByTokenAsync(token string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := ForgotUser{}
		err := self.db.Get(&user, "SELECT * from users_forgot_password WHERE token = ?", token)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) GetUserByIdAsync(id string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := User{}
		err := self.db.Get(&user, "SELECT * from users WHERE id = ?", id)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) UpdateAsync(user User) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		_, err := self.db.NamedExec("UPDATE users SET first_name=:first_name, last_name=:last_name, email=:email, username=:username, password=:password, confirm_token=:confirm_token, verified=:verified, reset=:reset WHERE id=:id", &user)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) DeleteForgotUser(id string) StoreResult {
	if id == "" {
		return StoreResult{
			Data:  nil,
			Total: 0,
			Err:   fmt.Errorf("null id"),
		}
	}

	_, err := self.db.Exec("DELETE from users_forgot_password where id = ?", id)

	return StoreResult{
		Data:  nil,
		Total: 0,
		Err:   err,
	}
}

func (self *UserSQLStore) SaveForgotUser(user ForgotUser) StoreResult {
	_, err := self.db.NamedExec("INSERT INTO users_forgot_password (id, token, created) values (:id, :token, :created)", &user)
	return StoreResult{
		Data:  user,
		Total: 1,
		Err:   err,
	}
}

func (self *UserSQLStore) Save(user User) StoreResult {
	_, err := self.db.Exec("INSERT INTO users "+
		"(id, first_name, last_name, email, username, password, confirm_token, verified, reset) "+
		"values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, user.FirstName, user.LastName, user.Email, user.Username, user.Password, user.ConfirmToken, false, false)
	return StoreResult{
		Data:  user,
		Total: 1,
		Err:   err,
	}
}
