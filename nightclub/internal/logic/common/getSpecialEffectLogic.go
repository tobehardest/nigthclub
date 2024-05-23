package common

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

var specialEffectNameMap = map[string]string{
	"1": "wechat_nightclub1_CZY_meigui.svga",
	"2": "wechat_nightclub1_CZY_meiyuan.svga",
	"3": "wechat_nightclub1_CZY_baoshijie.svga",
	"4": "wechat_nightclub1_CZY_chengbao.svga",
}

type GetSpecialEffectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	writer http.ResponseWriter
}

func NewGetSpecialEffectLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer http.ResponseWriter) *GetSpecialEffectLogic {
	return &GetSpecialEffectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		writer: writer,
	}
}

func (l *GetSpecialEffectLogic) GetSpecialEffect(req *types.GetSpecialEffectReq) error {
	// todo: add your logic here and delete this line
	//time.Sleep(60 * time.Second)
	//fileName := specialEffectMap[req.SpecialEffectId]
	//filePath := path.Join(constpath.VfxDir, fileName)
	//logx.Infof("download %s", filePath)
	//body, err := ioutil.ReadFile(filePath)
	//logx.Errorf("Read file length is %d", len(body))
	//if err != nil {
	//	return err
	//}
	body := l.svcCtx.ConstMap.SpecialEffectMap[req.SpecialEffectId]
	fileName := specialEffectNameMap[req.SpecialEffectId]
	l.writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)
	l.writer.Header().Add("Content-Description", "File Transfer")
	l.writer.Header().Add("Content-Transfer-Encoding", "binary")
	l.writer.Header().Add("Expires", "0")
	l.writer.Header().Add("Cache-Control", "must-revalidate")
	l.writer.Header().Add("Pragma", "public")
	n, err := l.writer.Write(body)
	logx.Errorf("The length of writer Write is %d", n)

	if err != nil {
		return err
	}

	if n < len(body) {
		return io.ErrClosedPipe
	}

	return nil
}
