package banip

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type IPRecord struct {
	IP           string
	EarliestTime time.Time //unix time
	LastTime     time.Time //unix time
	Count        int
	isIPBanned   bool
}

type FilterConfig struct {
	Skipper             middleware.Skipper
	ErrorHandler        func(context echo.Context, err error) error
	DenyHandler         func(context echo.Context, identifier string, err error) error
	IdentifierExtractor middleware.Extractor

	LimitRequestCount int
	LimitTime         time.Duration // second
	BanTime           time.Duration // second
	AllowList         []string
	DenyList          []string
}

var RecordMap sync.Map

func IsIPBanned(IP string, config FilterConfig) bool {
	if len(config.AllowList) > 0 {
		// if allow ip list is not empty, check if the ip is in the list
		for _, allowIP := range config.AllowList {
			if allowIP == IP {
				return false
			}
		}
	}

	// if deny ip list is not empty, check if the ip is in the list
	if len(config.DenyList) > 0 {
		for _, denyIP := range config.DenyList {
			if denyIP == IP {
				return true
			}
		}
	}

	if record, ok := RecordMap.Load(IP); ok {
		record := record.(IPRecord)
		// 判断是否已经被禁止
		if record.isIPBanned {
			// 判断laseTime是否超过限制时间
			if time.Now().Sub(record.LastTime) > config.BanTime {
				// 允许访问
				// 重置记录，从ban中拿出来
				RecordMap.Store(IP, IPRecord{
					IP:           IP,
					EarliestTime: time.Now(),
					LastTime:     time.Now(),
					Count:        1,
					isIPBanned:   false,
				})
				return false
			} else {
				return true
			}
		}

		// 如果当前时间减去最早访问时间大于限制时间，则重置记录
		if time.Now().Sub(record.EarliestTime) > config.LimitTime {
			RecordMap.Store(IP, IPRecord{
				IP:           IP,
				EarliestTime: time.Now(),
				LastTime:     time.Now(),
				Count:        1,
				isIPBanned:   false,
			})
			return false
		}

		// 否则判断是否超过限制次数
		if record.Count >= config.LimitRequestCount {
			// 超过限制次数，禁止访问
			RecordMap.Store(IP, IPRecord{
				IP:           IP,
				EarliestTime: record.EarliestTime,
				LastTime:     time.Now(),
				Count:        record.Count,
				isIPBanned:   true,
			})
			return true
		} else {
			// 正常访问的情况
			RecordMap.Store(IP, IPRecord{
				IP:           IP,
				EarliestTime: record.EarliestTime,
				LastTime:     time.Now(),
				Count:        record.Count + 1,
			})
			return false
		}

	} else {
		// 新访问情况
		RecordMap.Store(IP, IPRecord{
			IP:           IP,
			EarliestTime: time.Now(),
			LastTime:     time.Now(),
			Count:        1,
			isIPBanned:   false,
		})
		return false
	}
}
func FilterRequestConfig(config FilterConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = middleware.DefaultRateLimiterConfig.ErrorHandler
	}

	if config.DenyHandler == nil {
		config.DenyHandler = middleware.DefaultRateLimiterConfig.DenyHandler
	}

	if config.IdentifierExtractor == nil {
		config.IdentifierExtractor = middleware.DefaultRateLimiterConfig.IdentifierExtractor
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			identifier, err := config.IdentifierExtractor(c)

			if err != nil {
				c.Error(config.ErrorHandler(c, err))
			}

			// check if the ip is banned
			if IsIPBanned(identifier, config) {
				c.Error(config.DenyHandler(c, identifier, err))
				return nil
			}

			return next(c)
		}
	}
}
