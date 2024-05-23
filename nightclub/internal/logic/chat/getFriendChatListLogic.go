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

type GetFriendChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendChatListLogic {
	return &GetFriendChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendChatListLogic) GetFriendChatList(req *types.GetFriendChatListReq) (resp *types.GetFriendChatListResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	getHasNewFriendChatLogic := NewGetHasNewFriendChatLogic(l.ctx, l.svcCtx)
	hasNewFriendChat, err := getHasNewFriendChatLogic.GetHasNewFriendChat(&types.GetHasNewFriendChatReq{
		FromId: req.FromId,
		ToId:   userId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	if !hasNewFriendChat.HasNew {
		return &types.GetFriendChatListResp{}, nil
	}

	// 有新消息的话，获取一下索引位置和消息长度
	fromId := strconv.FormatInt(req.FromId, 10)
	toId := strconv.FormatInt(userId, 10)
	chatListHistoryKey := globalkey.ChatListHistoryKey + fromId + toId
	val, err := l.svcCtx.RedisClient.GetCtx(l.ctx, chatListHistoryKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	index, err := strconv.Atoi(val)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get conv string to int"), "fail to get conv string to int, err: %v", err)
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
	// 取索引位置到最后一条消息位置
	chatFriendList, err := l.svcCtx.RedisClient.LrangeCtx(l.ctx, chatFriendKey, index, len-1)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}

	getFriendChatListResp := &types.GetFriendChatListResp{}
	for _, c := range chatFriendList {
		chat := &types.Chat{}
		err := json.Unmarshal([]byte(c), chat)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
		}

		getFriendChatListResp.ChatList = append(getFriendChatListResp.ChatList, types.GetFriendChatVO{
			FromId:    chat.FromId,
			Content:   chat.Content,
			CreatTime: chat.CreatTime,
		})
	}

	for i, _ := range getFriendChatListResp.ChatList {
		fromId := getFriendChatListResp.ChatList[i].FromId
		id := strconv.FormatInt(fromId, 10)
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
		getFriendChatListResp.ChatList[i].UserName = user.UserName
		getFriendChatListResp.ChatList[i].Avatar = user.Avatar
	}

	// 存储history index
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, globalkey.ChatListHistoryKey+fromId+toId, strconv.Itoa(len))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	return getFriendChatListResp, nil
}
