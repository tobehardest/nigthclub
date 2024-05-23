package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
)

type UpdateFeatureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFeatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFeatureLogic {
	return &UpdateFeatureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFeatureLogic) UpdateFeature(req *types.UpdateFeatureReq) (resp *types.UpdateFeatureResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
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
	user.Features = req.Features
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userInfoKey marshal fail!"), "Load user data fail: %v", err)
	}
	l.svcCtx.RedisClient.SetCtx(l.ctx, userInfoKey, string(userJson))
	return &types.UpdateFeatureResp{
		Status: "success",
	}, nil
}
