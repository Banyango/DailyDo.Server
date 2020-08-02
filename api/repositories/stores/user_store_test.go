package stores

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"

)

func TestUserSQLStore_Save (t *testing.T) {
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

	userById := <- store.GetUserByIdAsync(id)

	assert.Nil(t, userById.Err)
	assert.NotNil(t, userById.Data)
}

func TestUserSQLStore_Update (t *testing.T) {
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

	<- store.UpdateAsync(updateUser)

	userById := <- store.GetUserByIdAsync(id)

	assert.Equal(t, updateUser, userById.Data.(model.User))
}
