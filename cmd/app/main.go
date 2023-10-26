package main

import (
	"github.com/RomanUtolin/RESTful-CRUD/internall/http"
	"github.com/RomanUtolin/RESTful-CRUD/internall/http/middleware"
	_logic "github.com/RomanUtolin/RESTful-CRUD/internall/logic"
	_repository "github.com/RomanUtolin/RESTful-CRUD/internall/repository"
	"github.com/RomanUtolin/RESTful-CRUD/pkg/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.GetLogger()
	ctx := config.GetConfigDb()
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
