package chat

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"
	"time"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendChatToFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendChatToFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendChatToFriendLogic {
	return &SendChatToFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendChatToFriendLogic) SendChatToFriend(req *types.SendFriendChatReq) (resp *types.SendFriendChatResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	fromId := strconv.FormatInt(userId, 10)
	toId := strconv.FormatInt(req.ToId, 10)
	var chatFriendKey string
	if fromId < toId {
		chatFriendKey = globalkey.ChatFriendKey + fromId + toId
	} else {
		chatFriendKey = globalkey.ChatFriendKey + toId + fromId
	}
	now := time.Now()
	chat := &types.Chat{
		ChatId:    now.Unix(),
		FromId:    userId,
		ToId:      req.ToId,
		Content:   req.Content,
		Status:    1,
		CreatTime: now.Format(time.DateTime),
	}

	chatBytes, err := json.Marshal(chat)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}

	_, err = l.svcCtx.RedisClient.RpushCtx(l.ctx, chatFriendKey, string(chatBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	// 存入聊天列表
	userChatListKey := globalkey.UserChatListKey + toId
	_, err = l.svcCtx.RedisClient.ZaddCtx(l.ctx, userChatListKey, now.Unix(), fromId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}
	return &types.SendFriendChatResp{
		Status: "success",
	}, nil
}
