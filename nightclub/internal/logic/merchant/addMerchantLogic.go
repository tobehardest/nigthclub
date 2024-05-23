package merchant

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"
	"time"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMerchantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMerchantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMerchantLogic {
	return &AddMerchantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMerchantLogic) AddMerchant(req *types.AddMerchantReq) (resp *types.AddMerchantResp, err error) {
	// todo: add your logic here and delete this line
	// 1.生成商户唯一id
	node, err := snowflake.NewNode(1)
	if err != nil {
		logx.WithContext(l.ctx).Error(err)
		return nil, err
	}
	id := node.Generate().Int64()
	merchantId := strconv.FormatInt(id, 10)
	logx.WithContext(l.ctx).Error("register userId id, %v", id)
	logx.WithContext(l.ctx).Error("register userInfo userId, %v", merchantId)
	now := time.Now()

	merchant := &types.Merchant{
		MerchantId:   id,
		MerchantName: req.MerchantName,
		Location:     req.Location,
		Status:       1,
		CreatTime:    now.Format(time.DateTime),
	}

	merChantBytes, err := json.Marshal(merchant)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}

	merchantKey := globalkey.MerchantKey + merchantId
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, merchantKey, string(merChantBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}
	return &types.AddMerchantResp{
		Status: "success",
	}, nil
}
