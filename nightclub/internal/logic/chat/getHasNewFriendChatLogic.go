package chat

import (
	"context"
	"github.com/pkg/errors"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHasNewFriendChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHasNewFriendChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHasNewFriendChatLogic {
	return &GetHasNewFriendChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHasNewFriendChatLogic) GetHasNewFriendChat(req *types.GetHasNewFriendChatReq) (resp *types.GetHasNewFriendChatResp, err error) {
	// todo: add your logic here and delete this line
	fromId := strconv.FormatInt(req.FromId, 10)
	toId := strconv.FormatInt(req.ToId, 10)
	chatListHistoryKey := globalkey.ChatListHistoryKey + fromId + toId
	exist, err := l.svcCtx.RedisClient.ExistsCtx(l.ctx, chatListHistoryKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	var chatFriendKey string
	if fromId < toId {
		chatFriendKey = globalkey.ChatFriendKey + fromId + toId
	} else {
		chatFriendKey = globalkey.ChatFriendKey + toId + fromId
	}
	len, err := l.svcCtx.RedisClient.LlenCtx(l.ctx, chatFriendKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	var index int
	if !exist {
		if len > 0 {
			err := l.svcCtx.RedisClient.SetCtx(l.ctx, chatListHistoryKey, "0")
			if err != nil {
				return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
			}
			return &types.GetHasNewFriendChatResp{
				HasNew: true,
			}, nil
		} else {
			return &types.GetHasNewFriendChatResp{
				HasNew: false,
			}, nil
		}
	} else {
		val, err := l.svcCtx.RedisClient.GetCtx(l.ctx, chatListHistoryKey)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
		}
		index, err = strconv.Atoi(val)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get conv string to int"), "fail to get conv string to int, err: %v", err)
		}
	}

	var flag bool
	if index <= len-1 {
		flag = true
	}
	return &types.GetHasNewFriendChatResp{
		HasNew: flag,
	}, nil
}
