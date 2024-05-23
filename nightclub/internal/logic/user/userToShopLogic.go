package user

import (
	"context"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserToShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserToShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserToShopLogic {
	return &UserToShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserToShopLogic) UserToShop(req *types.UserToStoreReq) (resp *types.UserToStoreResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	// 2、存入在线用户名单
	id := strconv.FormatInt(req.MerchantId, 10)
	onlineUserKey := globalkey.OnlineUserKey + id
	_, err = l.svcCtx.RedisClient.SaddCtx(l.ctx, onlineUserKey, userId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}
	return &types.UserToStoreResp{
		Status: "success",
	}, nil
}
