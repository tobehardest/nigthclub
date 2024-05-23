package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"nightclub/common/consturl"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/model"
	"strconv"
	"time"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdPaymentwxQueryPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentwxQueryPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentwxQueryPayLogic {
	return &ThirdPaymentwxQueryPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentwxQueryPayLogic) ThirdPaymentwxQueryPay(req *types.PaymentWxPayQueryReq) (resp *types.PaymentWxPayQueryResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	id := strconv.FormatInt(userId, 10)

	// get TransactionId
	thirdPayment := new(model.ThirdPayment)
	thirdPaymentOrderKey := globalkey.ThirdPaymentOrderKey + id + ":" + req.Sn
	thirdPaymentString, err := l.svcCtx.RedisClient.GetCtx(l.ctx, thirdPaymentOrderKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_SET_VALUE_ERROR, "fail to set value in redis"), "fail to set in redis, %v", err)
	}
	err = json.Unmarshal([]byte(thirdPaymentString), thirdPayment)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}

	ThirdPaymentWxPayQueryReq := &types.ThirdPaymentWxPayQueryReq{
		Platform:      "suixingfu",
		PayChannel:    "WX_LITE",
		SubClass:      "CONSUME",
		TenancyId:     l.svcCtx.Config.WxPayConf.MchId,
		StoreId:       0, // wait
		TransactionId: thirdPayment.TransactionId,
		OutTradeNo:    req.Sn,
	}
	body, _ := json.Marshal(ThirdPaymentWxPayQueryReq)
	request, err := http.NewRequest("POST", consturl.ThirdPaymentwxPayQueryUrl, bytes.NewBuffer(body))
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
	thirdPaymentWxPayQueryResp := &types.ThirdPaymentWxPayQueryResp{}
	json.Unmarshal(body, thirdPaymentWxPayQueryResp)

	return &types.PaymentWxPayQueryResp{}, nil

}
