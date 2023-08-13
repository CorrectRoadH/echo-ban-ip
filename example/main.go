package main

import (
	"net/http"
	"time"

	banip "github.com/CorrectRoadH/echo-ban-ip"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(banip.FilterRequestConfig(banip.FilterConfig{
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/favicon.ico" {
				return true
			}
			return false
		},
		LimitTime:         1 * time.Minute,
		LimitRequestCount: 60,
		BanTime:           1 * time.Hour,
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.String(http.StatusForbidden, "You are banned")
		},
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.Request().UserAgent(), nil
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
