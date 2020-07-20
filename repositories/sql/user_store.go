package sql

import (
	"fmt"
	"github.com/Banyango/gifoody_server/model"
	. "github.com/Banyango/gifoody_server/repositories/util"
	"github.com/google/uuid"
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


func (self *UserSQLStore) GetUserByEmail(email string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []model.Post
		rows, err := self.db.Query(`SELECT p.Id, p.name, p.url, p.user_id, p.post_date  FROM posts p LIMIT ? OFFSET ?`, query.Limit, query.Offset);
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}
		defer rows.Close()

		for rows.Next() {
			post := model.Post{}
			err := rows.Scan(&post.ID, &post.Name, &post.URL, &post.UserID, &post.PostDate)
			if err != nil {
				storeChan <- StoreResult{Data: nil, Err: err}
				return
			}
			results = append(results, post)
		}

		var count int
		row := self.db.QueryRow("SELECT COUNT(*) FROM posts")
		err = row.Scan(&count)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}

		storeChan <- StoreResult{Data: results, Total:count, Err: nil}
	}()
	return storeChan
}

func (self *UserSQLStore) GetUserByConfirmToken(token string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := []model.User{}
		err := self.db.Get(user, "SELECT * from user WHERE confirm_token = $1", token)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) GetUserById(id string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := []model.User{}
		err := self.db.Get(user, "SELECT * from user WHERE id = $1", id)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) Update(user model.User) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		tx := self.db.MustBegin()
		tx.MustExec("UPDATE user SET first_name=:first_name, last_name=:last_name, email=:email, username=:username, password=:password, confirm_token=:confirm_token, verified=:verified, reset=:reset WHERE id=:id", &user)
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

func (self *UserSQLStore) SaveForgotUser(user model.ForgotUser) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		tx := self.db.MustBegin()
		tx.MustExec("INSERT INTO user_forgot_password (id, token, created) values (:id, :token, :created)", &user)
		err := tx.Commit()
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *UserSQLStore) Save(user model.User) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		tx := self.db.MustBegin()
		user.Id = uuid.New().String()
		tx.MustExec("INSERT INTO user (id, first_name, last_name, email, username, password, confirm_token, verified, reset) values (:id, :first_name, :last_name, :email, :username, :password, :confirm_token, :verified, :reset)", &user)
		err := tx.Commit()
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}


