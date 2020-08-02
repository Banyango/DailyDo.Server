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


func (self *UserSQLStore) GetUserByEmailAsync(email string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := User{}
		err := self.db.Get(&user, "SELECT * from user WHERE email = ?", email)
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
		err := self.db.Get(&user, "SELECT * from user WHERE confirm_token = ?", token)
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
		err := self.db.Get(&user, "SELECT * from user WHERE id = ?", id)
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
		tx := self.db.MustBegin()
		tx.NamedExec("UPDATE user SET first_name=:first_name, last_name=:last_name, email=:email, username=:username, password=:password, confirm_token=:confirm_token, verified=:verified, reset=:reset WHERE id=:id", &user)
		err := tx.Commit()
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

	tx := self.db.MustBegin()
	tx.MustExec("DELETE from user where id = $1", id)
	err := tx.Commit()

	return StoreResult{
		Data:  nil,
		Total: 0,
		Err:   err,
	}
}

func (self *UserSQLStore) SaveForgotUser(user ForgotUser) StoreResult {
	tx := self.db.MustBegin()
	tx.MustExec("INSERT INTO user_forgot_password (id, token, created) values (:id, :token, :created)", &user)
	err := tx.Commit()
	return StoreResult{
		Data:  user,
		Total: 1,
		Err:   err,
	}
}

func (self *UserSQLStore) Save(user User) StoreResult {
	tx := self.db.MustBegin()
	tx.Query("INSERT INTO user (id, first_name, last_name, email, username, password, confirm_token, verified, reset) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Id, user.FirstName, user.LastName, user.Email, user.Username, user.Password, user.ConfirmToken, false, false)
	err := tx.Commit()
	return StoreResult{
		Data:  user,
		Total: 1,
		Err:   err,
	}
}


