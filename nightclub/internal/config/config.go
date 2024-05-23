package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	WxMiniConf WxMiniConf
	WxPayConf  WxPayConf
	Redis      redis.RedisConf
}
