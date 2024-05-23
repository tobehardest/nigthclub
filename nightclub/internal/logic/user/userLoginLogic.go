package user

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
	"strconv"
	"strings"
	"time"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	// 1、Get userId by mobile
	userId, err := l.loginByMobile(req.Mobile)
	if err != nil && userId == -1 {
		return nil, err
	}

	if userId == 0 {
		return nil, xerr.NewErrCodeMsg(10001, "用户未注册")
	}

	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	accessToken, err := generateTokenLogic.GenerateToken(&types.GenerateTokenReq{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_GENERATE_TOKEN_ERROR, "GenerateToken fail"), "GenerateToken fail userId: %d, err: %v", userId, err)
		//return nil, fmt.Errorf("GenerateToken fail userId : %d", userId)
	}
	_ = copier.Copy(resp, accessToken)

	// 修改用户最后登录时间
	id := strconv.FormatInt(userId, 10)
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
	now := time.Now()
	user.LastLoginTime = now.Unix()
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_FEATURES_UPDATE_ERROR, "userInfoKey marshal fail!"), "Load user data fail: %v", err)
	}
	l.svcCtx.RedisClient.SetCtx(l.ctx, userInfoKey, string(userJson))

	loginResp := &types.LoginResp{
		UserId:      userId,
		AccessToken: *accessToken,
	}
	return loginResp, nil
}

func (l *UserLoginLogic) loginByMobile(mobile string) (int64, error) {
	mobileKey := globalkey.UserMobileKey + mobile
	logx.WithContext(l.ctx).Info("hcy 手机号，mobile:%s", mobile)
	val, err := l.svcCtx.RedisClient.Get(mobileKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		logx.WithContext(l.ctx).Errorf("根据手机号查询用户信息失败，mobile:%s,err:%v", mobile, err)
		return -1, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "根据手机号查询用户信息失败"), "mobile:%s,err:%v", mobile, err)
		//return -1, fmt.Errorf("根据手机号查询用户信息失败，mobile:%s,err:%w", mobile, err)
	}
	if val == "" {
		logx.WithContext(l.ctx).Errorf("用户不存在, mobile:%s,err:%w", mobile, err)
		return 0, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_NOT_FINDINFO_ERROR, "用户不存在"), "mobile:%s,err:%v", mobile, err)
		//return 0, fmt.Errorf("用户不存在, mobile:%s,err:%w", mobile, err)
	}
	logx.WithContext(l.ctx).Info("hcy val，val:%s", val)

	userKey := strings.Replace(val, globalkey.UserInfoKey, "", 1)
	logx.WithContext(l.ctx).Info("hcy userKey，val:%s", userKey)
	userId, _ := strconv.ParseInt(userKey, 10, 64)
	logx.WithContext(l.ctx).Info("hcy userId，userId:%s", userId)
	return userId, nil
}
