package user

import (
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
	"time"
)

type UserRegisterLogic struct {
	logx.Logger
	r      *http.Request
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(r *http.Request, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(r.Context()),
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	mobileKey := globalkey.UserMobileKey + req.Mobile
	userId, err := l.svcCtx.RedisClient.Get(mobileKey)
	if err != nil && err != redis.Nil {
		logx.WithContext(l.r.Context()).Errorf("user already register")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_ALREADY_REGIST_ERROR, "user already register"), "user already register, err: %v", err)
	}

	// 号码已经注册了
	if userId != "" {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_ALREADY_REGIST_ERROR, "user already register"), "user already register, err: %v", err)
	}

	// 1.生成用户唯一id
	node, err := snowflake.NewNode(1)
	if err != nil {
		logx.WithContext(l.r.Context()).Error(err)
		return nil, err
	}
	id := node.Generate().Int64()
	userId = strconv.FormatInt(id, 10)
	logx.WithContext(l.r.Context()).Error("register userId id, %v", id)
	logx.WithContext(l.r.Context()).Error("register userInfo userId, %v", userId)

	// 2.存储用户信息
	now := time.Now()
	user := &types.User{
		UserId:        id,
		Mobile:        req.Mobile,
		LastLoginTime: now.Unix(),
		CreatTime:     now.Format(time.DateTime),
	}
	userInfoKey := globalkey.UserInfoKey + userId
	userBytes, err := json.Marshal(user)
	if err != nil {
		logx.WithContext(l.r.Context()).Error("fail to json marshal, %v", err)
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "fail to json marshal"), "fail to json marshal, %v", err)
	}
	err = l.svcCtx.RedisClient.SetCtx(l.r.Context(), userInfoKey, string(userBytes))
	if err != nil {
		logx.WithContext(l.r.Context()).Error("fail to set userInfoKey, %v", err)
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_SETINFO_ERROR, "Sava userInfo fail"), "fail to set userInfoKey, %v", err)
	}

	logx.WithContext(l.r.Context()).Error("register userInfoKey, %v", userInfoKey)
	err = l.svcCtx.RedisClient.SetCtx(l.r.Context(), mobileKey, userInfoKey)
	if err != nil {
		logx.WithContext(l.r.Context()).Error("fail to set mobileKey, %v", err)
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_SETINFO_ERROR, "fail to set mobile"), "fail to set mobileKey, err: %v", err)
	}

	//2、Generate the token, so that the service doesn't call rpc internally
	generateTokenLogic := NewGenerateTokenLogic(l.r.Context(), l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&types.GenerateTokenReq{
		UserId: id,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.SERVER_COMMON_ERROR, "GenerateToken userId fail"), "GenerateToken userId fail, err: %v", err)
	}

	logx.WithContext(l.r.Context()).Errorf("regMsg userId :", userId)
	registerResp := &types.RegisterResp{
		UserId:      id,
		AccessToken: *tokenResp,
	}

	return registerResp, nil
}
