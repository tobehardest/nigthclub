package order

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/utils"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/model"
	"strconv"
	"time"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatSendGiftOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatSendGiftOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatSendGiftOrderLogic {
	return &CreatSendGiftOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatSendGiftOrderLogic) CreatSendGiftOrder(req *types.CreatSendGiftOrderReq) (resp *types.CreatSendGiftOrderResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	fromId := strconv.FormatInt(userId, 10)
	now := time.Now()

	// 1.生成订单唯一id
	node, err := snowflake.NewNode(1)
	if err != nil {
		logx.WithContext(l.ctx).Error(err)
		return nil, err
	}
	id := node.Generate().Int64()
	// sendGiftOrderId := strconv.FormatInt(id, 10)
	logx.WithContext(l.ctx).Error("register userId id, %v", id)
	logx.WithContext(l.ctx).Error("register userInfo userId, %v", userId)

	// 2.存入订单
	sendGiftOrderModel := &model.SendGiftOrder{
		Id:         id,
		Sn:         utils.GenSn(utils.SN_PREFIX_HOMESTAY_ORDER),
		FromId:     userId,
		ToId:       req.ToId,
		GiftPrice:  req.GiftPrice,
		TradeState: model.SendGiftOrderTradeStateWaitPay,
		TradeCode:  utils.Krand(8, utils.KC_RAND_KIND_ALL),
		CreatTime:  now,
		UpdateTime: now,
	}
	sendGiftOrderModelBytes, err := json.Marshal(sendGiftOrderModel)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}
	sendGiftOrderKey := globalkey.SendGiftOrderKey + fromId + ":" + sendGiftOrderModel.Sn
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, sendGiftOrderKey, string(sendGiftOrderModelBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	return &types.CreatSendGiftOrderResp{
		OrderSn: sendGiftOrderModel.Sn,
	}, nil
}
