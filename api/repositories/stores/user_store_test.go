package stores

import (
	"github.com/Banyango/dailydo_server/api/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserSQLStore_Save(t *testing.T) {
	db, err := sqlx.Connect("mysql", "fooduser:foodtest@/food_test?parseTime=true")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	store := NewUserSQLStore(db)

	id := uuid.New().String()
	user := model.User{
		Id:           id,
		Email:        uuid.New().String(),
		FirstName:    uuid.New().String(),
		LastName:     uuid.New().String(),
		Password:     uuid.New().String(),
		Username:     uuid.New().String(),
		ConfirmToken: uuid.New().String(),
		Reset:        false,
		Verified:     false,
	}

	store.Save(user)

	userById := <-store.GetUserByIdAsync(id)

	assert.Nil(t, userById.Err)
	assert.NotNil(t, userById.Data)
}

func TestUserSQLStore_SaveShouldFail_WhenDuplicateEmail(t *testing.T) {
	db, err := sqlx.Connect("mysql", "fooduser:foodtest@/food_test?parseTime=true")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	store := NewUserSQLStore(db)

	id := uuid.New().String()
	user := model.User{
		Id:           id,
		Email:        "alreadyUsed@email.com",
		FirstName:    uuid.New().String(),
		LastName:     uuid.New().String(),
		Password:     uuid.New().String(),
		Username:     uuid.New().String(),
		ConfirmToken: uuid.New().String(),
		Reset:        false,
		Verified:     false,
	}

	user2 := model.User{
		Id:           uuid.New().String(),
		Email:        "alreadyUsed@email.com",
		FirstName:    uuid.New().String(),
		LastName:     uuid.New().String(),
		Password:     uuid.New().String(),
		Username:     uuid.New().String(),
		ConfirmToken: uuid.New().String(),
		Reset:        false,
		Verified:     false,
	}

	store.Save(user)
	result := store.Save(user2)

	assert.NotNil(t, result.Err)
}

func TestUserSQLStore_SaveShouldFail_WhenDuplicateUsername(t *testing.T) {
	db, err := sqlx.Connect("mysql", "fooduser:foodtest@/food_test?parseTime=true")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	store := NewUserSQLStore(db)

	id := uuid.New().String()
	user := model.User{
		Id:           id,
		Email:        uuid.New().String(),
		FirstName:    uuid.New().String(),
		LastName:     uuid.New().String(),
		Password:     uuid.New().String(),
		Username:     "user1",
		ConfirmToken: uuid.New().String(),
		Reset:        false,
		Verified:     false,
	}

	user2 := model.User{
		Id:           uuid.New().String(),
		Email:        uuid.New().String(),
		FirstName:    uuid.New().String(),
		LastName:     uuid.New().String(),
		Password:     uuid.New().String(),
		Username:     "user1",
		ConfirmToken: uuid.New().String(),
		Reset:        false,
		Verified:     false,
	}

	store.Save(user)
	result := store.Save(user2)

	assert.NotNil(t, result.Err)
}

func TestUserSQLStore_Update(t *testing.T) {
	db, err := sqlx.Connect("mysql", "fooduser:foodtest@/food_test?parseTime=true")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	store := NewUserSQLStore(db)

	id := uuid.New().String()
	originalUser := model.User{
		Id:           id,
		Email:        uuid.New().String(),
		FirstName:    uuid.New().String(),
		LastName:     uuid.New().String(),
		Password:     uuid.New().String(),
		Username:     uuid.New().String(),
		ConfirmToken: uuid.New().String(),
		Reset:        false,
		Verified:     false,
	}

	store.Save(originalUser)

	updateUser := model.User{
		Id:           id,
		Email:        "fart@fart.com",
		FirstName:    "John",
		LastName:     "lolcopters",
		Password:     "123",
		Username:     "tty",
		ConfirmToken: "123",
		Reset:        true,
		Verified:     true,
	}

	<-store.UpdateAsync(updateUser)

	userById := <-store.GetUserByIdAsync(id)

	assert.Equal(t, updateUser, userById.Data.(model.User))
}
