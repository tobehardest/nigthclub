package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWechatIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWechatIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWechatIdLogic {
	return &UpdateWechatIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWechatIdLogic) UpdateWechatId(req *types.UpdateWechatIdReq) (resp *types.UpdateWechatIdResp, err error) {
	// todo: add your logic here and delete this line
	// todo: 这个UserId 应该是客户端发来的
	//userId := ctxdata.GetUidFromCtx(l.ctx)
	id := strconv.FormatInt(req.UserId, 10)
	userInfoKey := globalkey.UserInfoKey + id
	val, err := l.svcCtx.RedisClient.Get(userInfoKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Connect RedisClient fail!"), "redis search fail, err: %v", err)
	}
	user := new(types.User)
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("Unmarshal fail!")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Unmarshal json fromat fail!"), "Unmarshal fail! err: %v", err)
	}
	user.WechatId = req.WechatId
	userJson, err := json.Marshal(user)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("user Marshal fail!")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "userJson marshal userId fail!"), "user Marshal fail! err: %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, userInfoKey, string(userJson))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("user Marshal fail!")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "userInfoKey marshal fail!"), "user Marshal fail! err: %v", err)
	}
	return &types.UpdateWechatIdResp{
		Status: "success",
	}, nil
}
