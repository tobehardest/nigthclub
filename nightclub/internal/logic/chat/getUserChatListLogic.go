package chat

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

type GetUserChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserChatListLogic {
	return &GetUserChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserChatListLogic) GetUserChatList(req *types.GetUserChatListReq) (resp *types.GetUserChatListResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	toId := strconv.FormatInt(userId, 10)
	userChatListKey := globalkey.UserChatListKey + toId
	start := (req.Page - 1) * req.Size
	end := start + req.Size - 1
	userChatList, err := l.svcCtx.RedisClient.ZrevrangeCtx(l.ctx, userChatListKey, start, end)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	getUserChatListResp := &types.GetUserChatListResp{}
	for _, fromId := range userChatList {
		getUserChatVO := &types.GetUserChatVO{}
		// get user info
		userInfokey := globalkey.UserInfoKey + fromId
		userString, err := l.svcCtx.RedisClient.GetCtx(l.ctx, userInfokey)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
		}
		user := new(types.User)
		err = json.Unmarshal([]byte(userString), user)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Unmarshal json fromat fail!"), "Unmarshal fail! err: %v", err)
		}
		getUserChatVO.FromId = user.UserId
		getUserChatVO.UserName = user.UserName
		getUserChatVO.Avatar = user.Avatar
		getUserChatVO.Features = user.Features

		// get Last message
		var chatFriendKey string
		if fromId < toId {
			chatFriendKey = globalkey.ChatFriendKey + fromId + toId
		} else {
			chatFriendKey = globalkey.ChatFriendKey + toId + fromId
		}
		lastChatMessage, err := l.svcCtx.RedisClient.LindexCtx(l.ctx, chatFriendKey, -1)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
		}
		chat := &types.Chat{}
		err = json.Unmarshal([]byte(lastChatMessage), chat)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
		}
		getUserChatVO.Content = chat.Content
		getUserChatVO.CreatTime = chat.CreatTime
		getUserChatListResp.UserChatList = append(getUserChatListResp.UserChatList, *getUserChatVO)
	}
	return getUserChatListResp, nil
}
