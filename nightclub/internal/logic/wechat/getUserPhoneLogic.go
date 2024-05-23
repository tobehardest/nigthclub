package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"time"
)

type GetUserPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPhoneLogic {
	return &GetUserPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var (
	mp = make(map[string]types.WxAccessToken)
)

func (l *GetUserPhoneLogic) GetUserPhone(req *types.GetPhoneReq) (resp *types.GetPhoneResp, err error) {
	accessToken, err := GetAccessToken(l.ctx, req.Code)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}

	reqBody := types.GetPhoneReq{
		Code: req.Code,
	}
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}
	logx.WithContext(l.ctx).Infof("request gpt json string : %v", string(reqData))
	// TODO 没有正确拿到手机号，只拿到一个空字符串
	request, err := http.NewRequest("POST", globalkey.WechatGetUserPhoneNumberUrl+accessToken, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		logx.WithContext(l.ctx).Errorf("getPhoneNumber api status code not equals 200,code is %d ,details:  %v ", res.StatusCode, string(body))
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("failed to read response body")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}
	logx.WithContext(l.ctx).Infof("response gpt json string : %v", string(body))

	phoneBack := new(types.WxPhoneBack)
	err = json.Unmarshal(body, phoneBack)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("failed to unmarshal the response body")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}

	if phoneBack.ErrCode == -1 {
		logx.WithContext(l.ctx).Errorf("The system is busy, you need to resend the request")
		res, err = client.Do(request)
		body, err = ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, phoneBack)
		if phoneBack.ErrCode == -1 {
			return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
		}
	}
	if phoneBack.ErrCode == 40029 {
		logx.WithContext(l.ctx).Errorf("The system is busy, you need to resend the request")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}

	var reply string
	reply = phoneBack.PhoneInfo.PurePhoneNumber
	logx.WithContext(l.ctx).Infof("gpt response text: %s ", reply)
	ans := &types.GetPhoneResp{
		PurePhoneNumber: reply,
	}
	return ans, nil
}

func GetAccessToken(ctx context.Context, code string) (string, error) {
	// 先查本地缓存
	tmp := types.WxAccessToken{}
	if accessTokenRes := mp[globalkey.RequestIdKey]; accessTokenRes != tmp {
		if int64(time.Now().Unix()) < accessTokenRes.ExpiresIn {
			logx.Info("get access token from local cache")
			return accessTokenRes.AccessToken, nil
		}
	}
	// 本地缓存无效，请求服务器
	logx.Info("Get access token from wc server")
	req, err := http.NewRequest("GET", globalkey.WechatAccessTokenUrl, nil)
	if err != nil {
		return "", errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(xerr.NewErrCodeMsg(xerr.ACCESS_WECHATAPI_ERROR, "fail to access wechat API"), "fail to access wechat API, err: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		logx.WithContext(ctx).Errorf("getPhoneNumber api status code not equals 200,code is %d ,details:  %v ", res.StatusCode, string(body))
		return "", errors.New(fmt.Sprintf("getPhoneNumber api status code not equals 200,code is %d ,details:  %v ", res.StatusCode, string(body)))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logx.WithContext(ctx).Errorf("failed to read response body")
		return "", err
	}
	logx.WithContext(ctx).Infof("response gpt json string : %v", string(body))

	accessTokenRes := new(types.WxAccessToken)
	err = json.Unmarshal(body, accessTokenRes)
	if err != nil {
		logx.WithContext(ctx).Errorf("failed to unmarshal the response body")
		return "", err
	}

	// 存入本地缓存
	logx.WithContext(ctx).Infof("Store the access token in the local cache")
	accessTokenRes.ExpiresIn += int64(time.Now().Unix()) - 1200
	mp[globalkey.RequestIdKey] = *accessTokenRes
	return accessTokenRes.AccessToken, nil
}
