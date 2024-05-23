package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLocationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLocationLogic {
	return &UpdateUserLocationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLocationLogic) UpdateUserLocation(req *types.UpdateUserLocationReq) (resp *types.UpdateUserLocationResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	id := strconv.FormatInt(userId, 10)
	userInfoKey := globalkey.UserInfoKey + id
	val, err := l.svcCtx.RedisClient.GetCtx(l.ctx, userInfoKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "Connect RedisClient fail!"), "Redis query fail err: %v", err)
	}
	user := new(types.User)
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}
	user.Location = req.Location
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userInfoKey marshal fail!"), "Load user data fail: %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, userInfoKey, string(userJson))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	return &types.UpdateUserLocationResp{
		Status: "success",
	}, nil
}
