package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func init() {
	viper.SetConfigFile("configs/server.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfigDb() string {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	sslMode := viper.GetString(`database.sslmode`)
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode)
	return dbUrl
}

func GetLogger() {
	logLevel, _ := logrus.ParseLevel(viper.GetString("log_level"))
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if viper.GetBool("debug") {
		logrus.Infof("Service RUN on DEBUG mode")
	}
}

func GetDb() *pgxpool.Pool {
	ctx := context.Background()
	dbUrl := GetConfigDb()
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		logrus.Warning(err)
	}
	return dbPool
}

func GetTimeoutContext() time.Duration {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	return timeoutContext
}
