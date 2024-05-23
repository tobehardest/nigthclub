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

type SendPatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendPatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPatMessageLogic {
	return &SendPatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPatMessageLogic) SendPatMessage(req *types.SendPatMessageReq) (resp *types.SendPatMessageResp, err error) {
	// todo: add your logic here and delete this line
	now := time.Now()
	message := &types.Message{
		MessageId: now.Unix(),
		FromId:    req.FromId,
		ToId:      req.ToId,
		Type:      globalkey.TOPIC_PAT,
		Status:    1,
		CreatTime: now.Format(time.DateTime),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}

	toId := strconv.FormatInt(req.ToId, 10)
	messageKey := globalkey.MessageKey + toId
	_, err = l.svcCtx.RedisClient.ZaddCtx(l.ctx, messageKey, now.Unix(), string(messageBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30001, "redis 查询失败"), "fail search redis, err: %v", err)
	}

	return &types.SendPatMessageResp{
		Status: "success",
	}, nil
}
