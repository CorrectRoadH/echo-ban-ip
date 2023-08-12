package banip

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type IPRecord struct {
	IP           string
	EarliestTime time.Time //unix time
	LastTime     time.Time //unix time
	Count        int
	isIPBanned   bool
}

type FilterConfig struct {
	LimitRequestCount int
	LimitTime         time.Duration // second
	BanTime           time.Duration // second
	AllowIPList       []string
	DenyIPList        []string
	DenyPrompt        string
}

func getIPFromRealIP(ip string) string {
	return ip[:strings.LastIndex(ip, ":")]
}

var RecordMap map[string]IPRecord = make(map[string]IPRecord)

func IsIPBanned(IP string, config FilterConfig) bool {
	if len(config.AllowIPList) > 0 {
		// if allow ip list is not empty, check if the ip is in the list
		for _, allowIP := range config.AllowIPList {
			if allowIP == IP {
				return false
			}
		}
	}

	// if deny ip list is not empty, check if the ip is in the list
	if len(config.DenyIPList) > 0 {
		for _, denyIP := range config.DenyIPList {
			if denyIP == IP {
				return true
			}
		}
	}

	if record, ok := RecordMap[IP]; ok {
		// 判断是否已经被禁止
		if record.isIPBanned {
			// 判断laseTime是否超过限制时间
			if time.Now().Sub(record.LastTime) > config.BanTime {
				// 允许访问
				// 重置记录，从ban中拿出来
				RecordMap[IP] = IPRecord{
					IP:           IP,
					EarliestTime: time.Now(),
					LastTime:     time.Now(),
					Count:        1,
					isIPBanned:   false,
				}
				return false
			} else {
				return true
			}
		}

		// 如果当前时间减去最早访问时间大于限制时间，则重置记录
		if time.Now().Sub(record.EarliestTime) > config.LimitTime {
			RecordMap[IP] = IPRecord{
				IP:           IP,
				EarliestTime: time.Now(),
				LastTime:     time.Now(),
				Count:        1,
				isIPBanned:   false,
			}
			return false
		}

		// 否则判断是否超过限制次数
		if record.Count >= config.LimitRequestCount {
			// 超过限制次数，禁止访问
			RecordMap[IP] = IPRecord{
				IP:           IP,
				EarliestTime: record.EarliestTime,
				LastTime:     time.Now(),
				Count:        record.Count,
				isIPBanned:   true,
			}
			return true
		} else {
			// 正常访问的情况
			RecordMap[IP] = IPRecord{
				IP:           IP,
				EarliestTime: record.EarliestTime,
				LastTime:     time.Now(),
				Count:        record.Count + 1,
			}
			return false
		}

	} else {
		// 新访问情况
		RecordMap[IP] = IPRecord{
			IP:           IP,
			EarliestTime: time.Now(),
			LastTime:     time.Now(),
			Count:        1,
			isIPBanned:   false,
		}
		return false
	}
}
func FilterRequestConfig(config FilterConfig) echo.MiddlewareFunc {
	if config.DenyPrompt == "" {
		config.DenyPrompt = "Your IP is banned cause of too many requests in a short time."
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get the ip address
			IP := getIPFromRealIP(c.RealIP())

			// check if the ip is banned
			if IsIPBanned(IP, config) {
				return c.String(http.StatusForbidden, config.DenyPrompt)
			}

			return next(c)
		}
	}
}
