package http

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func RequestsLogger() echo.MiddlewareFunc {
	return logger
}

func logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return logRequest(c, next)
	}
}

func logRequest(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()

	start := time.Now()

	if err := next(c); err != nil {
		c.Error(err)
	}

	latency := time.Since(start).Milliseconds()

	log.WithFields(log.Fields{
		"remote_ip":  c.RealIP(),
		"host":       req.Host,
		"uri":        req.RequestURI,
		"method":     req.Method,
		"user_agent": req.UserAgent(),
		"status":     res.Status,
		"latency":    latency,
	}).Debug("request processed")

	return nil
}

func RequestsBodiesLogger() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.WithFields(log.Fields{
			"request_body":  string(reqBody),
			"response_body": string(resBody),
		}).Debug("request bodies processed")
	})
}
