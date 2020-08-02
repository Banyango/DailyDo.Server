package users
//
//import (
//	"encoding/json"
//	"github.com/Banyango/gifoody_server/api/infrastructure/mail"
//	"github.com/Banyango/gifoody_server/api/infrastructure/os"
//	"github.com/Banyango/gifoody_server/api/infrastructure/template"
//	"github.com/Banyango/gifoody_server/api/repositories"
//	"github.com/jmoiron/sqlx"
//	"github.com/labstack/echo/v4"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//func setup_test_store() *repositories.AppStore {
//	db, _ := sqlx.Connect("mysql", "fooduser:foodtest@/food_test?parseTime=true")
//	return repositories.NewAppStore(db)
//}
//
//func TestUserController_PostRegister_ShouldSaveUser(t *testing.T) {
//
//	if testing.Short() {
//		t.Skip()
//	}
//
//	e := echo.New()
//	registerJson := CreateUserRequest{
//		Email:     "abc@123.com",
//		Password:  "123",
//		Username:  "user1",
//		FirstName: "kylyg",
//		LastName:  "recwas",
//	}
//	marshal, _ := json.Marshal(registerJson)
//	userJSON := strings.NewReader(string(marshal))
//	req := httptest.NewRequest(http.MethodPost, "/register", userJSON)
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	store := setup_test_store()
//	h := NewUserController(store.User(), , )
//
//	// Assertions
//	if assert.NoError(t, h.PostRegister(c)) {
//		assert.Equal(t, http.StatusCreated, rec.Code)
//		assert.Equal(t, userJSON, rec.Body.String())
//	}
//
//	var response map[string]string
//	json.Unmarshal([]byte(rec.Body.String()), &response)
//
//	user := <- store.User().GetUserByIdAsync(response["id"])
//
//	assert.NotNil(t, user.Data)
//}
//
