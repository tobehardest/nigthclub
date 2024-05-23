package payment

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/model"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrWxPayCallbackError = xerr.NewErrMsg("wechat pay callback fail")

type ThirdPaymentWxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentWxPayCallbackLogic {
	return &ThirdPaymentWxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentWxPayCallbackLogic) ThirdPaymentWxPayCallback(req *types.ThirdPaymentWxPayCallbackReq) (resp *types.ThirdPaymentWxPayCallbackResp, err error) {
	// todo: add your logic here and delete this line

	// todo : 验证签名

	// Verify and update relevant flow data
	userId := ctxdata.GetUidFromCtx(l.ctx)
	fromId := strconv.FormatInt(userId, 10)
	thirdPaymentOrderKey := globalkey.ThirdPaymentOrderKey + fromId + ":" + req.OrdNo
	thirdPaymentOrderString, err := l.svcCtx.RedisClient.GetCtx(l.ctx, thirdPaymentOrderKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	thirdPayment := new(model.ThirdPayment)
	err = json.Unmarshal([]byte(thirdPaymentOrderString), thirdPayment)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}
	// 对比金额
	totalOffstAmt := req.TotalOffstAmt
	notifyPayTotal, err := strconv.ParseInt(totalOffstAmt, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(ErrWxPayCallbackError, "Failed to get payment flow record err:%v ,notifyTrasaction:%+v ", err, req)
	}
	if thirdPayment.PayTotal != notifyPayTotal {
		return nil, errors.Wrapf(ErrWxPayCallbackError, "Failed to get payment flow record err:%v ,notifyTrasaction:%+v ", err, req)
	}

	// Judgment status
	if req.BizCode == "success" {
		//Payment Notification.

		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return nil, errors.Wrapf(ErrWxPayCallbackError, "Failed to get payment flow record err:%v ,notifyTrasaction:%+v ", err, req)
		}

		// Update the flow status.
		if _, err = l.UpdateTradeState(&types.UpdateTradeStateReq{
			Sn:             req.OrdNo,
			TradeState:     "SUCCESS",
			TransactionId:  req.TransactionId,
			TradeType:      req.ActivityNo,
			TradeStateDesc: req.Detail,
			PayStatus:      model.ThirdPaymentPayTradeStateSuccess,
		}); err != nil {
			return nil, errors.Wrapf(ErrWxPayCallbackError, "更新流水状态失败  err:%v , notifyTrasaction:%v ", err, req)
		}

	} else if req.BizCode == "paying" {
		//Refund notification @todo to be done later, not needed at this time
	}

	return &types.ThirdPaymentWxPayCallbackResp{
		Code: "success",
		Msg:  "成功",
	}, nil
}

func (l *ThirdPaymentWxPayCallbackLogic) UpdateTradeState(req *types.UpdateTradeStateReq) (resp *types.UpdateTradeStateResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	fromId := strconv.FormatInt(userId, 10)
	//1、payment record confirm
	thirdPaymentOrderKey := globalkey.ThirdPaymentOrderKey + fromId + ":" + req.Sn
	thirdPaymentOrderString, err := l.svcCtx.RedisClient.GetCtx(l.ctx, thirdPaymentOrderKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	thirdPayment := new(model.ThirdPayment)
	err = json.Unmarshal([]byte(thirdPaymentOrderString), thirdPayment)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}

	if thirdPayment == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("third payment record no exists"), " sn : %s", req.Sn)
	}

	//2、Judgment Status
	if req.PayStatus == model.ThirdPaymentPayTradeStateSuccess || req.PayStatus == model.ThirdPaymentPayTradeStateFAIL {
		//Want to modify as payment success, failure scenarios
		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return &types.UpdateTradeStateResp{}, nil
			//return nil, errors.Wrapf(xerr.NewErrMsg("Only orders with wait payment can be success of fail"), "Only orders with waitting payment can be success for fail in : %+v", req)
		}

	} else if req.PayStatus == model.ThirdPaymentPayTradeStateRefund {
		//Want to change to refund success scenario

		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateSuccess {
			return nil, errors.Wrapf(xerr.NewErrMsg("Only orders with successful payment can be refunded"), "Only orders with successful payment can be refunded in : %+v", req)
		}
	} else {
		return nil, errors.Wrapf(xerr.NewErrMsg("This status is not currently supported"), "Modify payment flow status is not supported  in : %+v", req)
	}

	//3、update .
	thirdPayment.TradeState = req.TradeState
	thirdPayment.TransactionId = req.TransactionId
	thirdPayment.TradeType = req.TradeType
	thirdPayment.TradeStateDesc = req.TradeStateDesc
	thirdPayment.PayStatus = req.PayStatus
	thirdPayment.PayTime = time.Unix(req.PayTime, 0)
	thirdPayment.UpdateTime = time.Now()
	thirdPaymentBytes, err := json.Marshal(thirdPayment)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, thirdPaymentOrderKey, string(thirdPaymentBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	//4、notify  sub "payment-update-paystatus-topic"  services(order-mq ..), pub、sub use kq

	return &types.UpdateTradeStateResp{}, nil
}
