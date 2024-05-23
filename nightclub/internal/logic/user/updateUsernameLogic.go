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

type UpdateUserNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserNameLogic {
	return &UpdateUserNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserNameLogic) UpdateUserName(req *types.UpdateUserNameReq) (resp *types.UpdateUserNameResp, err error) {
	// todo: add your logic here and delete this line
	//userId := ctxdata.GetUidFromCtx(l.ctx)
	userInfoKey := globalkey.UserInfoKey + strconv.FormatInt(req.UserId, 10)
	// userInfoKey := globalkey.UserInfoKey + req.UserId
	val, err := l.svcCtx.RedisClient.Get(userInfoKey)
	if err != nil {
		//logx.WithContext(l.ctx).Errorf("userInfoKey get redisClient fail!")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Connect RedisClient fail!"), "userInfoKey get redisClient fail, err: %v", err)
	}
	user := new(types.User)
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		// logx.WithContext(l.ctx).Errorf("Unmarshal fail!")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Unmarshal json fromat fail!"), "Unmarshal fail! err: %v", err)
	}

	user.UserName = req.UserName
	userJson, err := json.Marshal(user)
	if err != nil {
		// logx.WithContext(l.ctx).Errorf("user Marshal fail!")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "userJson marshal userId fail!"), "user Marshal fail! err: %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, userInfoKey, string(userJson))
	if err != nil {
		//logx.WithContext(l.ctx).Errorf("user Marshal fail!")\
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "userInfoKey marshal fail!"), "user set fail,err: %v", err)
	}
	return &types.UpdateUserNameResp{
		Status: "success",
	}, nil
}
