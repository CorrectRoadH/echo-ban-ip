# What is echo-ban-ip? 

`echo-ban-ip` is a very easy to use echo middleware that to ban ip address that request too many. to prevent simple script attack.


# Usage
    
**Basic**
```go
import (

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
    ban_ip "github.com/CorrectRoadH/go-echo-ban-ip"
)
package main
func main(){
    e := echo.New()
    // if request count > 100 in 10 minute, ban ip for 1 hour
    e.Use(ban_ip.FilterRequestConfig(ban_ip.FilterConfig{
		LimitRequestCount: 100,
		LimitTime:         time.Minute * 10,
		BanTime:           time.Hour * 1,
	}))

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}
```

**AllowList and DenyList**
```go
import (

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
    ban_ip "github.com/CorrectRoadH/go-echo-ban-ip"
)
package main
func main(){
    e := echo.New()
    // if request count > 100 in 10 minute, ban ip for 1 hour
    e.Use(ban_ip.FilterRequestConfig(ban_ip.FilterConfig{
		LimitRequestCount: 100,
		LimitTime:         time.Minute * 10,
		BanTime:           time.Hour * 1,
        AllowIPList:       []string{"127.0.0.1"},
        DenyIPList:        []string{"192.168.5.2"},
	}))

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}
```


**custom deny prompt**
```go
import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
    ban_ip "github.com/CorrectRoadH/go-echo-ban-ip"
)
package main
func main(){
    e := echo.New()
    // if request count > 100 in 10 minute, ban ip for 1 hour
    e.Use(ban_ip.FilterRequestConfig(ban_ip.FilterConfig{
		LimitRequestCount: 100,
		LimitTime:         time.Minute * 10,
		BanTime:           time.Hour * 1,
        DenyPrompt:        "Access Deny",
	}))

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}
```
