package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlineUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnlineUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineUsersLogic {
	return &OnlineUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnlineUsersLogic) OnlineUsers(req *types.OnlineUsersReq) (resp *types.OnlineUsersResp, err error) {
	// todo: add your logic here and delete this line
	currentId := ctxdata.GetUidFromCtx(l.ctx)

	// get merchantId
	id := strconv.FormatInt(req.MerchantId, 10)
	onlineUserKey := globalkey.OnlineUserKey + id
	userIdList, err := l.svcCtx.RedisClient.SmembersCtx(l.ctx, onlineUserKey)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("OnlineUsersResp get redisClient fail!")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	userInfoPrefix := globalkey.UserInfoKey

	users := []types.User{}
	for _, userId := range userIdList {
		// 跳过当前登录用户
		nowId, err := strconv.ParseInt(userId, 10, 64)
		if nowId == currentId {
			continue
		}
		userInfoKey := userInfoPrefix + userId
		userString, err := l.svcCtx.RedisClient.GetCtx(l.ctx, userInfoKey)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
		}
		user := new(types.User)
		err = json.Unmarshal([]byte(userString), user)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
		}

		// 检查的登录是否已过期
		accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
		now := time.Now()
		if user.LastLoginTime+accessExpire < now.Unix() {
			_, err := l.svcCtx.RedisClient.SremCtx(l.ctx, onlineUserKey, userId)
			if err != nil {
				return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
			}
			continue
		}
		users = append(users, *user)
	}

	return &types.OnlineUsersResp{
		List: users,
	}, nil
}
