package middleware

import (
	"context"
	"github.com/RomanUtolin/RESTful-CRUD/internall/constants"
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type GoMiddleware struct {
	ctx context.Context
}

func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

func (m *GoMiddleware) LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		c.Response().After(func() {
			logrus.WithFields(logrus.Fields{
				"method":     c.Request().Method,
				"path":       c.Path(),
				"code":       c.Response().Status,
				"latency_ns": time.Now().Sub(start).Nanoseconds(),
			}).Info("request Details")
		})
		return next(c)
	}
}

func (m *GoMiddleware) DbSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dbSession := newDbSession(m.ctx)
		c.Set(constants.DbSession, dbSession)
		c.Response().After(func() {
			err := dbSession.Close()
			if err != nil {
				logrus.Warning(err)
			}
		})
		return next(c)
	}
}

func InitMiddleware(ctx context.Context) *GoMiddleware {
	return &GoMiddleware{ctx: ctx}
}

func newDbSession(ctx context.Context) *dbr.Session {
	dbDriver := ctx.Value(constants.DatabaseDriver).(string)
	dbUrl := ctx.Value(constants.DatabaseURL).(string)
	dbConn, err := dbr.Open(dbDriver, dbUrl, nil)
	if err != nil {
		logrus.Warning(err)
	}
	err = dbConn.Ping()
	if err != nil {
		logrus.Warning(err)
	}
	dbSession := dbConn.NewSession(nil)
	return dbSession
}
