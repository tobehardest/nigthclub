package message

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

type GetHasSendGiftMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHasSendGiftMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHasSendGiftMessageLogic {
	return &GetHasSendGiftMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHasSendGiftMessageLogic) GetHasSendGiftMessage(req *types.GetHasSendGiftReq) (resp *types.GetHasSendGiftResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	fromId := strconv.FormatInt(userId, 10)
	toId := strconv.FormatInt(req.ToId, 10)

	// 查对方有没有给自己送过礼物
	hasSendGiftKey := globalkey.HasSendGiftKey + toId
	hasSendGift, err := l.svcCtx.RedisClient.SismemberCtx(l.ctx, hasSendGiftKey, fromId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	if hasSendGift {
		return &types.GetHasSendGiftResp{
			HasSendGift: true,
		}, nil
	}

	// 查自己有没有给对方送给礼物
	hasSendGiftKey = globalkey.HasSendGiftKey + fromId
	hasSendGift, err = l.svcCtx.RedisClient.SismemberCtx(l.ctx, hasSendGiftKey, toId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	if hasSendGift {
		return &types.GetHasSendGiftResp{
			HasSendGift: true,
		}, nil
	}

	return &types.GetHasSendGiftResp{
		HasSendGift: false,
	}, nil
}
