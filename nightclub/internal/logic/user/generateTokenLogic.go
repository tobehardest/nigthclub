package user

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
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

type GenerateTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateTokenLogic) GenerateToken(req *types.GenerateTokenReq) (resp *types.GenerateTokenResp, err error) {
	// todo: add your logic here and delete this line
	now := time.Now()
	// id, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
	}
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now.Unix(), accessExpire, req.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.TOKEN_GENERATE_ERROR, "generate token fail"), "getJwtToken err userId:%d , err:%v", req.UserId, err)
	}

	// 修改用户最后登录时间
	id := strconv.FormatInt(req.UserId, 10)
	userInfoKey := globalkey.UserInfoKey + id
	val, err := l.svcCtx.RedisClient.Get(userInfoKey)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "Connect RedisClient fail!"), "Redis query fail err: %v", err)
	}
	user := new(types.User)
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userJson marshal userId fail!"), "format conversion fail: %v", err)
	}
	user.LastLoginTime = now.Unix()
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userInfoKey marshal fail!"), "Load user data fail: %v", err)
	}
	l.svcCtx.RedisClient.SetCtx(l.ctx, userInfoKey, string(userJson))

	return &types.GenerateTokenResp{
		AccessToken:  accessToken,
		AccessExpire: now.Unix() + accessExpire,
		RefreshAfter: now.Unix() + accessExpire/2,
	}, nil
}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {

	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[ctxdata.CtxKeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
