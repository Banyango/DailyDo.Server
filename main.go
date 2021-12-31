package main

import (
	"github.com/Banyango/dailydo_server/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	. "net/http"
	"os"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.AddTrailingSlash())

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	db, err := sqlx.Connect("mysql", dbConnectionString)

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(3 * time.Second)

	defer db.Close()

	_, isDevExists := os.LookupEnv("DEV")

	config := middleware.CORSConfig{
		AllowOrigins: []string{"https://192.168.1.79"},
		AllowMethods: []string{MethodGet, MethodPut, MethodPost, MethodDelete, MethodOptions},
	}

	if isDevExists {
		config.AllowOrigins = []string{"http://localhost:3000"}
	}

	e.Use(middleware.CORSWithConfig(config))

	api.InitRouter(e, db)

	//e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":3001"))

}
