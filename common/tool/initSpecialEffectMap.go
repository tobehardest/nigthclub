package tool

import (
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"nightclub/common/constpath"
	"path"
)

//var specialEffectMap = map[string]string{
//	"1": "wechat_nightclub1_CZY_meigui.svga",
//	"2": "wechat_nightclub1_CZY_meiyuan.svga",
//	"3": "wechat_nightclub1_CZY_baoshijie.svga",
//	"4": "wechat_nightclub1_CZY_chengbao.svga",
//}

func initSpecialEffectMap() map[string][]byte {
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
	return specialEffectMap
}
