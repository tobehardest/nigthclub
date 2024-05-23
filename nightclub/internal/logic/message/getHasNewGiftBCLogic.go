package message

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHasNewGiftBCLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type HasNewGiftBC struct {
	FromName string `json:"from_name"`
	ToName   string `json:"to_name"`
	GiftId   int64  `json:"giftId"`
	Time     int64  `json:"time"`
}

func NewGetHasNewGiftBCLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHasNewGiftBCLogic {
	return &GetHasNewGiftBCLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHasNewGiftBCLogic) GetHasNewGiftBC(req *types.GetHasNewGiftBCReq) (resp *types.GetHasNewGiftBCResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	id := strconv.FormatInt(userId, 10)
	MerchantId := strconv.FormatInt(req.MerchantId, 10)
	giftBroadCastHistoryKey := globalkey.GiftBroadCastHistory + MerchantId + ":" + id
	giftBroadCastKey := globalkey.GiftBroadCastKey + MerchantId
	exist, err := l.svcCtx.RedisClient.ExistsCtx(l.ctx, giftBroadCastHistoryKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	llen, err := l.svcCtx.RedisClient.Llen(giftBroadCastKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	var index int
	if !exist {
		if llen > 0 {
			err := l.svcCtx.RedisClient.SetCtx(l.ctx, giftBroadCastHistoryKey, "0")
			if err != nil {
				return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
			}
			index = 0
		} else {
			return &types.GetHasNewGiftBCResp{
				HasNew: false,
			}, nil
		}
	} else {
		giftBroadCastHistory, err := l.svcCtx.RedisClient.GetCtx(l.ctx, giftBroadCastHistoryKey)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
		}
		index, err = strconv.Atoi(giftBroadCastHistory)
	}

	if index > llen-1 {
		return &types.GetHasNewGiftBCResp{
			HasNew: false,
		}, nil
	}

	hasNewGiftBCList, err := l.svcCtx.RedisClient.LrangeCtx(l.ctx, giftBroadCastKey, index, llen-1)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	getHasNewGiftBCResp := new(types.GetHasNewGiftBCResp)
	hasNewGiftBC := new(HasNewGiftBC)
	now := time.Now().Unix()
	for _, hasNewGiftBCString := range hasNewGiftBCList {
		err := json.Unmarshal([]byte(hasNewGiftBCString), hasNewGiftBC)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Marshal fail"), "fail Marshal val, err: %v", err)
		}
		if hasNewGiftBC.Time+3600 < now {
			continue
		}
		getHasNewGiftBCResp.GetHasNewGiftBCList = append(getHasNewGiftBCResp.GetHasNewGiftBCList, types.GetHasNewGiftBCVO{
			GiftId:   hasNewGiftBC.GiftId,
			FromName: hasNewGiftBC.FromName,
			ToName:   hasNewGiftBC.ToName,
		})
	}

	// update index
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, giftBroadCastHistoryKey, strconv.Itoa(llen))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	getHasNewGiftBCResp.HasNew = true
	return getHasNewGiftBCResp, nil
}
