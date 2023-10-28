package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type GoMiddleware struct {
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

func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
