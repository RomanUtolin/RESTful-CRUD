package config

import (
	"context"
	"fmt"
	"github.com/RomanUtolin/RESTful-CRUD/internall/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("configs/server.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfigDb() context.Context {
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
	return ctx
}

func GetLogger() {
	logLevel, _ := logrus.ParseLevel(viper.GetString("log_level"))
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if viper.GetBool("debug") {
		logrus.Infof("Service RUN on DEBUG mode")
	}
}
