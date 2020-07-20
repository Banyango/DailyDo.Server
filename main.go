package main

import (
	"github.com/Banyango/gifoody_server/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	. "net/http"
)

func main() {

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.AddTrailingSlash())

	db, err := sqlx.Connect("mysql", "fooduser:foodtest@/food_test?parseTime=true")
	defer db.Close()

	if err != nil {
		panic(err)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://api.gifoody.com", "http://localhost:3000"},
		AllowMethods: []string{MethodGet, MethodPut, MethodPost, MethodDelete, MethodOptions},
	}))

	api.InitRouter(e, db)

	//e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":3001"))

}
