package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"io/ioutil"
	"net/http"
	"nightclub/common/consturl"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/tool"
	"nightclub/common/utils"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/model"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrWxMiniAuthFailError = xerr.NewErrMsg("wechat mini auth fail")
var ErrWxPayError = xerr.NewErrMsg("wechat pay fail")

type ThirdPaymentwxPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentwxPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentwxPayLogic {
	return &ThirdPaymentwxPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentwxPayLogic) ThirdPaymentwxPay(req *types.PaymentWxPayReq) (*types.PaymentWxPayResp, error) {
	// todo: add your logic here and delete this line
	// 1、get user openId
	userId := ctxdata.GetUidFromCtx(l.ctx)
	id := strconv.FormatInt(userId, 10)
	now := time.Now()
	var totalPrice int64 // Total amount paid for current order(cent)
	var description string

	// 2、 get order description
	switch req.ServiceType {
	case model.ThirdPaymentServiceTypeSendGiftOrder:
		sendGiftTotalPrice, sendGiftDescription, err := l.getPaySendGiftPriceDesciption(req.OrderSn)
		if err != nil {
			return nil, errors.Wrapf(ErrWxPayError, "getPayHomestayPriceDescription err : %v req: %+v", err, req)
		}
		totalPrice, description = sendGiftTotalPrice, sendGiftDescription

	default:
		return nil, errors.Wrapf(xerr.NewErrMsg("Payment for this business type is not supported"), "Payment for this business type is not supported req: %+v", req)
	}

	// 3、get oepnId
	miniProgram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.Secret,
		Cache:     cache.NewMemory(),
	})
	authResult, err := miniProgram.GetAuth().Code2Session(req.Code)
	if err != nil || authResult.ErrCode != 0 || authResult.OpenID == "" {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "发起授权请求失败 err : %v , code : %s  , authResult : %+v", err, req.Code, authResult)
	}

	// 4、create local third payment record
	node, err := snowflake.NewNode(1)
	if err != nil {
		logx.WithContext(l.ctx).Error(err)
		return nil, err
	}
	thirdPaymentId := node.Generate().Int64()
	// sendGiftOrderId := strconv.FormatInt(id, 10)
	logx.WithContext(l.ctx).Error("register userId id, %v", id)
	logx.WithContext(l.ctx).Error("register userInfo userId, %v", userId)
	thirdPayment := new(model.ThirdPayment)
	thirdPayment.Id = thirdPaymentId
	thirdPayment.UserId = userId
	thirdPayment.PayMode = model.ThirdPaymentPayModelWechatPay
	thirdPayment.Sn = utils.GenSn(utils.SN_PREFIX_THIRD_PAYMENT)
	thirdPayment.PayTotal = totalPrice
	thirdPayment.OrderSn = req.OrderSn
	thirdPayment.ServiceType = model.ThirdPaymentServiceTypeSendGiftOrder
	thirdPayment.PayStatus = model.ThirdPaymentPayTradeStateWait
	thirdPayment.CreatTime = now
	thirdPayment.UpdateTime = now
	thirdPaymentBytes, err := json.Marshal(thirdPayment)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}
	thirdPaymentOrderKey := globalkey.ThirdPaymentOrderKey + id + ":" + thirdPayment.Sn
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, thirdPaymentOrderKey, string(thirdPaymentBytes))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}

	// 3、 create wechat pay pre pay order
	thirdPaymentWxPayReq := &types.ThirdPaymentWxPayReq{
		Platform:   "suixingfu",
		PayChannel: "WX_LITE",
		SubClass:   "CONSUME",
		TenancyId:  l.svcCtx.Config.WxPayConf.MchId,
		StoreId:    0,
		OutTradeNo: thirdPayment.Sn,
		Openid:     authResult.OpenID,
		Body:       description,
		Amount:     float64(totalPrice),
		NotifyUrl:  l.svcCtx.Config.WxPayConf.NotifyUrl,
	}
	body, _ := json.Marshal(thirdPaymentWxPayReq)
	url := tool.AddUrlSign(consturl.ThirdPaymentwxPayUrl, "1")
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.HTTP_CREAT_HTTP_ERROR, "failed to create post request"), "failed to create post request, err: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 100 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.HTTP_SEND_REQUEST_ERROR, "failed to create post request"), "failed to create post request, err: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ = ioutil.ReadAll(response.Body)
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.HTTP_STATUS_CODE_ERROR, "api status code not equals 200"), "api status code not equals 200,code is %d ,details:  %v ", response.StatusCode, string(body))
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.HTTP_READ_RES_ERROR, "failed to read response body"), "failed to read response body, err: %v", err)
	}
	thirdPaymentWxPayResp := &types.ThirdPaymentWxPayResp{}
	json.Unmarshal(body, thirdPaymentWxPayResp)

	return &types.PaymentWxPayResp{
		Appid:     thirdPaymentWxPayResp.JsApiPayConfig.AppId,
		NonceStr:  thirdPaymentWxPayResp.JsApiPayConfig.NonceStr,
		PaySign:   thirdPaymentWxPayResp.JsApiPayConfig.PaySign,
		Package:   thirdPaymentWxPayResp.JsApiPayConfig.Package,
		Timestamp: thirdPaymentWxPayResp.JsApiPayConfig.Timestamp,
		SignType:  thirdPaymentWxPayResp.JsApiPayConfig.SignType,
	}, nil
}

func (l *ThirdPaymentwxPayLogic) getPaySendGiftPriceDesciption(orderSn string) (int64, string, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	id := strconv.FormatInt(userId, 10)
	description := "send git pay"

	// get user price
	sendGiftOrderKey := globalkey.SendGiftOrderKey + id + ":" + orderSn
	sendGiftString, err := l.svcCtx.RedisClient.GetCtx(l.ctx, sendGiftOrderKey)
	if err != nil {
		return 0, description, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	sendGiftOrder := new(model.SendGiftOrder)
	err = json.Unmarshal([]byte(sendGiftString), sendGiftOrder)
	if err != nil {
		return 0, description, errors.Wrapf(ErrWxPayError,
			"OrderRpc.HomestayOrderDetail err: %v, orderSn: %s", err, orderSn)
	}

	if sendGiftOrder == nil || sendGiftOrder.Id == 0 {
		return 0, description, errors.Wrapf(xerr.NewErrMsg("order no exists"), "WeChat payment order does not exist orderSn : %s", orderSn)
	}

	return sendGiftOrder.GiftPrice, description, nil
}
