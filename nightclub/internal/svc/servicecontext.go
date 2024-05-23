package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"io/ioutil"
	"nightclub/common/constpath"
	"nightclub/nightclub/internal/config"
	"path"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis
	ConstMap    *ConstMap
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		RedisClient: redis.MustNewRedis(c.Redis),
		ConstMap:    InitConstMap(),
	}
}

type ConstMap struct {
	SpecialEffectMap map[string][]byte
}

func InitConstMap() *ConstMap {
	constMap := new(ConstMap)
	constMap.initSpecialEffectMap()
	return constMap
}

func (constMap *ConstMap) initSpecialEffectMap() {
	//fileName := specialEffectMap[req.SpecialEffectId]
	specialEffectMap := make(map[string][]byte)
	fileName := "wechat_nightclub1_CZY_meigui.svga"
	filePath := path.Join(constpath.VfxDir, fileName)
	logx.Infof("download %s", filePath)
	body, _ := ioutil.ReadFile(filePath)
	logx.Errorf("Read file length is %d", len(body))
	specialEffectMap["1"] = body

	fileName = "wechat_nightclub1_CZY_meiyuan.svga"
	filePath = path.Join(constpath.VfxDir, fileName)
	logx.Infof("download %s", filePath)
	body, _ = ioutil.ReadFile(filePath)
	logx.Errorf("Read file length is %d", len(body))
	specialEffectMap["2"] = body

	fileName = "wechat_nightclub1_CZY_baoshijie.svga"
	filePath = path.Join(constpath.VfxDir, fileName)
	logx.Infof("download %s", filePath)
	body, _ = ioutil.ReadFile(filePath)
	logx.Errorf("Read file length is %d", len(body))
	specialEffectMap["3"] = body

	fileName = "wechat_nightclub1_CZY_chengbao.svga"
	filePath = path.Join(constpath.VfxDir, fileName)
	logx.Infof("download %s", filePath)
	body, _ = ioutil.ReadFile(filePath)
	logx.Errorf("Read file length is %d", len(body))
	specialEffectMap["4"] = body
	constMap.SpecialEffectMap = specialEffectMap
}
