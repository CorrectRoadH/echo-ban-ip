package banip_test

import (
	"testing"
	"time"

	middleware "github.com/CorrectRoadH/echo-ban-ip"

	"github.com/stretchr/testify/assert"
)

func TestIsIPBannedCase1(t *testing.T) {
	config := middleware.FilterConfig{
		LimitTime:         time.Second * 1,
		LimitRequestCount: 2,
		BanTime:           time.Second * 5,
	}
	result := middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, false, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, false, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, true, result)
	time.Sleep(time.Second * 6)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, false, result)
}

func TestIsIPBannedAllowList(t *testing.T) {
	config := middleware.FilterConfig{
		LimitTime:         time.Second * 1,
		LimitRequestCount: 2,
		BanTime:           time.Second * 5,
		AllowList:         []string{"127.0.0.1"},
	}
	result := middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, false, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, false, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, false, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, false, result)
}

func TestIsIPBannedDenyList(t *testing.T) {
	config := middleware.FilterConfig{
		LimitTime:         time.Second * 1,
		LimitRequestCount: 2,
		BanTime:           time.Second * 5,
		DenyList:          []string{"127.0.0.1"},
	}
	result := middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, true, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, true, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, true, result)
	result = middleware.IsIPBanned("127.0.0.1", config)
	assert.Equal(t, true, result)
}
