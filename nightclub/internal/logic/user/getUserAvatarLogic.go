package user

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"nightclub/common/constpath"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"path"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	writer http.ResponseWriter
}

func NewGetUserAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer http.ResponseWriter) *GetUserAvatarLogic {
	return &GetUserAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		writer: writer,
	}
}

func (l *GetUserAvatarLogic) GetUserAvatar(req *types.GetUserAvatarReq) (err error) {
	filePath := path.Join(constpath.ImgDir, req.File)
	logx.Infof("download %s", filePath)
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	n, err := l.writer.Write(body)
	if err != nil {
		return err
	}

	if n < len(body) {
		return io.ErrClosedPipe
	}

	return nil
}
