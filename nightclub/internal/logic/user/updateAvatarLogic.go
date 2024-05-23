package user

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"nightclub/common/constpath"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"os"
	"path"
	"strconv"
)

const maxFileSize = 10 << 24

type UpdateAvatarLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUpdateAvatarLogic(r *http.Request, svcCtx *svc.ServiceContext) *UpdateAvatarLogic {
	return &UpdateAvatarLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAvatarLogic) UpdateAvatar(req *types.UpdateAvatarReq) (resp *types.UpdateAvatarResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.r.Context())
	l.r.ParseMultipartForm(maxFileSize)
	file, header, err := l.r.FormFile("avatar")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logx.Infof("upload file: %+v, file size: %d, MIME header: %+v",
		header.Filename, header.Size, header.Header)

	fileName := header.Filename
	//originFileName := header.Filename

	filePath := path.Join(constpath.ImgDir, fileName)
	tempFile, err := os.Create(filePath)

	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	io.Copy(tempFile, file)

	// 存入 userInfo
	id := strconv.FormatInt(userId, 10)
	userInfoKey := globalkey.UserInfoKey + id
	val, err := l.svcCtx.RedisClient.Get(userInfoKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "Connect RedisClient fail!"), "Redis query fail err: %v", err)
	}
	user := new(types.User)
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}
	user.Avatar = filePath
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userInfoKey marshal fail!"), "Load user data fail: %v", err)
	}
	l.svcCtx.RedisClient.SetCtx(l.r.Context(), userInfoKey, string(userJson))

	return &types.UpdateAvatarResp{
		Avatar: fileName,
	}, nil

}
