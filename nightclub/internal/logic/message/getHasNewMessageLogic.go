package message

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"
	"strings"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHasNewMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHasNewMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHasNewMessageLogic {
	return &GetHasNewMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHasNewMessageLogic) GetHasNewMessage(req *types.GetHasNewMessageReq) (resp *types.GetHasNewMessageResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	id := strconv.FormatInt(userId, 10)

	// 查询余额
	userInfo, err := l.svcCtx.RedisClient.GetCtx(l.ctx, globalkey.UserInfoKey+id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	user := new(types.User)
	err = json.Unmarshal([]byte(userInfo), user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.USER_UNSERNAME_UPDATE_ERROR, "Unmarshal json fromat fail!"), "Unmarshal fail! err: %v", err)
	}
	getHasNewMessageResp := &types.GetHasNewMessageResp{}
	getHasNewMessageResp.Balance = user.Balance

	// 获取是否有新消息
	messageNoticeHistory := globalkey.MessageNoticeHistory + id
	index, err := l.svcCtx.RedisClient.GetCtx(l.ctx, messageNoticeHistory)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	messageKey := globalkey.MessageKey + id
	last, err := l.svcCtx.RedisClient.ZrevrangeCtx(l.ctx, messageKey, 0, 0)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.Redis_GET_VAULE_ERROR, "fail to get from redis"), "fail to get from redis, err: %v", err)
	}
	if last == nil || len(last) == 0 {
		getHasNewMessageResp.HasNewMessage = false
		return getHasNewMessageResp, nil
	}
	message := &types.Message{}
	err = json.Unmarshal([]byte(last[0]), message)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(30002, "Unmarshal fail"), "fail Unmarshal val, err: %v", err)
	}

	if strings.Compare(index, message.CreatTime) == -1 {
		getHasNewMessageResp.HasNewMessage = true
	} else {
		getHasNewMessageResp.HasNewMessage = false
	}
	return getHasNewMessageResp, nil
}
