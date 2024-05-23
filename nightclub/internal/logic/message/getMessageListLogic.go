package message

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageListLogic) GetMessageList(req *types.GetMessageListReq) (resp *types.GetMessageListResp, err error) {
	// todo: add your logic here and delete this line
	id := strconv.FormatInt(req.ToId, 10)
	messageKey := globalkey.MessageKey + id
	start := (req.Page - 1) * req.Size
	end := start + req.Size - 1
	messages, err := l.svcCtx.RedisClient.Zrevrange(messageKey, start, end)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30001, "redis 查询失败"), "fail search redis, err: %v", err)
	}
	getMessageListResp := &types.GetMessageListResp{
		MessageList: []types.Message{},
	}
	for _, m := range messages {
		message := &types.Message{}
		err := json.Unmarshal([]byte(m), message)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
		}

		getMessageListResp.MessageList = append(getMessageListResp.MessageList, *message)
	}

	messageNoticeHistory := globalkey.MessageNoticeHistory + id
	last, err := l.svcCtx.RedisClient.ZrevrangeCtx(l.ctx, messageKey, 0, 0)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	message := &types.Message{}
	err = json.Unmarshal([]byte(last[0]), message)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, messageNoticeHistory, message.CreatTime)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}
	return getMessageListResp, nil
}
