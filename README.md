# What is echo-ban-ip? 

`echo-ban-ip` is a out of the box echo middleware that to ban ip address that request too many. to prevent simple script attack.
`echo-ban-ip` 是一个开箱即用的echo中间件，用于禁止请求过多的ip地址。以防止简单的脚本攻击。

# What is the difference between this and Rate Limiter?

The middleware is consistent of the most function of [Rate Limiter](https://echo.labstack.com/docs/middleware/rate-limiter). The most difference is that can ban ip address that request too many for punish.
该中间件与[Rate Limiter](https://echo.labstack.com/docs/middleware/rate-limiter)的大部分功能一致。最大的区别是可以禁止请求过多的ip地址一个额外时长以进行惩罚。
# Usage

**Basic**
```go
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
		LimitTime:         1 * time.Minute,
		LimitRequestCount: 60,
		BanTime:           1 * time.Hour,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
```

**AllowList and DenyList**
```go
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
		LimitTime:         1 * time.Minute,
		LimitRequestCount: 60,
		BanTime:           1 * time.Hour,
		AllowList:         []string{"127.0.0.1"},
		DenyList:         []string{"192.168.1.1"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

```


**Custom Deny Hander**
```go
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
		LimitTime:         1 * time.Minute,
		LimitRequestCount: 60,
		BanTime:           1 * time.Hour,
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.String(http.StatusForbidden, "You are banned")
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
```

**Custom Identifier**
```go
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
```

**Skipper**
```go
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
```