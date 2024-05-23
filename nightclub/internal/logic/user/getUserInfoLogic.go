package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.GetUserInfoResp, err error) {
	// todo: add your logic here and delete this line
	id := strconv.FormatInt(req.UserId, 10)
	userInfokey := globalkey.UserInfoKey + id
	val, err := l.svcCtx.RedisClient.GetCtx(l.ctx, userInfokey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logx.WithContext(l.ctx).Errorf("user data not find")
			return nil, xerr.NewErrCodeMsg(xerr.USER_NOT_FINDINFO_ERROR, "user not find")
		}
		logx.WithContext(l.ctx).Errorf("fail to find RedisClient data, %v", err)
		return nil, err
	}

	if val == "" {
		return nil, xerr.NewErrCodeMsg(xerr.USER_NOT_FINDINFO_ERROR, "user info is null")
	}

	user := new(types.User)
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("fail to unmarshal data, %v", err)
		return nil, err
	}

	return &types.GetUserInfoResp{
		User: *user,
	}, nil
}
