package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
)

type GetUserWechatIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserWechatIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWechatIdLogic {
	return &GetUserWechatIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserWechatIdLogic) GetUserWechatId(req *types.GetUserWechatIdReq) (resp *types.GetUserWechatIdResp, err error) {
	// todo: add your logic here and delete this line
	id := strconv.FormatInt(req.UserId, 10)
	userInfokey := globalkey.UserInfoKey + id
	val, err := l.svcCtx.RedisClient.GetCtx(l.ctx, userInfokey)
	if err != nil || val == "" {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_GET_WECHATID_ERROR, "UserKey is not exist err"), "User is not exist err: %v, user: %d", err, req.UserId)
	}
	user := new(types.User)
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_GET_WECHATID_ERROR, "Format conversion fail"), "Format conversion fail err: %v", err)
	}
	getUserWechatIdResp := &types.GetUserWechatIdResp{
		WechatId: user.WechatId,
	}
	return getUserWechatIdResp, nil
}
