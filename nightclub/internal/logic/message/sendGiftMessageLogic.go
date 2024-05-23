package message

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"
	"time"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendGiftMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendGiftMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendGiftMessageLogic {
	return &SendGiftMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendGiftMessageLogic) SendGiftMessage(req *types.SendGiftMessageReq) (resp *types.SendGiftMessageResp, err error) {
	// todo: add your logic here and delete this line
	fromId := strconv.FormatInt(req.FromId, 10)
	toId := strconv.FormatInt(req.ToId, 10)
	giftContent := req.GiftContent

	// 添加送过礼物
	hasSendGiftKey := globalkey.HasSendGiftKey + toId
	_, err = l.svcCtx.RedisClient.SaddCtx(l.ctx, hasSendGiftKey, fromId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	// 修改双方余额
	fromUserInfo, err := l.svcCtx.RedisClient.GetCtx(l.ctx, globalkey.UserInfoKey+fromId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	fromUser := new(types.User)
	err = json.Unmarshal([]byte(fromUserInfo), fromUser)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Unmarshal json fromat fail!"), "Unmarshal fail! err: %v", err)
	}
	fromUser.Balance -= giftContent.GiftAmount
	userJson, err := json.Marshal(fromUser)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userInfoKey marshal fail!"), "Load user data fail: %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, globalkey.UserInfoKey+fromId, string(userJson))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	toUserInfo, err := l.svcCtx.RedisClient.GetCtx(l.ctx, globalkey.UserInfoKey+toId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	toUser := new(types.User)
	err = json.Unmarshal([]byte(toUserInfo), toUser)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Unmarshal json fromat fail!"), "Unmarshal fail! err: %v", err)
	}
	toUser.Balance += giftContent.GiftAmount
	userJson, err = json.Marshal(toUser)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userInfoKey marshal fail!"), "Load user data fail: %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, globalkey.UserInfoKey+toId, string(userJson))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	// 发送通知
	now := time.Now()
	giftContentBytes, err := json.Marshal(giftContent)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}
	message := &types.Message{
		MessageId: now.Unix(),
		FromId:    req.FromId,
		ToId:      req.ToId,
		Type:      globalkey.TOPIC_Gift,
		Content:   string(giftContentBytes),
		Status:    1,
		CreatTime: now.Format(time.DateTime),
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}
	messageKey := globalkey.MessageKey + toId
	_, err = l.svcCtx.RedisClient.ZaddCtx(l.ctx, messageKey, now.Unix(), string(messageBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}

	// 发送广播
	hasNewGiftBC := &HasNewGiftBC{
		FromName: fromUser.UserName,
		ToName:   toUser.UserName,
		GiftId:   req.GiftContent.GiftId,
		Time:     now.Unix(),
	}
	hasNewGiftBCBytes, err := json.Marshal(hasNewGiftBC)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "fail to json marshal"), "fail to json marshal, %v", err)
	}
	giftBroadCastKey := globalkey.GiftBroadCastKey + strconv.FormatInt(req.MerchantId, 10)
	_, err = l.svcCtx.RedisClient.RpushCtx(l.ctx, giftBroadCastKey, string(hasNewGiftBCBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	return &types.SendGiftMessageResp{
		Status: "success",
	}, nil
}
