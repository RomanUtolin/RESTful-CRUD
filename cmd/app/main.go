package main

import (
	"context"
	"fmt"
	"github.com/RomanUtolin/RESTful-CRUD/internall/constants"
	"github.com/RomanUtolin/RESTful-CRUD/internall/http"
	"github.com/RomanUtolin/RESTful-CRUD/internall/http/middleware"
	_logic "github.com/RomanUtolin/RESTful-CRUD/internall/logic"
	_repository "github.com/RomanUtolin/RESTful-CRUD/internall/repository"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("configs/server.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	logLevel, _ := logrus.ParseLevel(viper.GetString("log_level"))
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if viper.GetBool("debug") {
		logrus.Infof("Service RUN on DEBUG mode")
	}
}

func main() {
	dbDriver := viper.GetString("database.driver")
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	sslMode := viper.GetString(`database.sslmode`)
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode)
	ctx := context.WithValue(context.Background(), constants.DatabaseDriver, dbDriver)
	ctx = context.WithValue(ctx, constants.DatabaseURL, dbUrl)

	server := echo.New()
	middl := middleware.InitMiddleware(ctx)
	server.Use(middl.CORS)
	server.Use(middl.LogRequest)
	server.Use(middl.DbSession)
	repository := _repository.NewPersonRepository()
	logic := _logic.NewPersonLogic(repository)
	http.NewHandler(server, logic)

	logrus.Infof("Starting Server")
	err := server.Start(viper.GetString("server.address"))
	if err != nil {
		logrus.Warning(err)
	}
}
