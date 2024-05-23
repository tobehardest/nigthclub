package chat

import (
	"context"
	"github.com/pkg/errors"
	"nightclub/common/ctxdata"
	"nightclub/common/globalkey"
	"nightclub/common/xerr"
	"strconv"

	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHasNewUserChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHasNewUserChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHasNewUserChatLogic {
	return &GetHasNewUserChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHasNewUserChatLogic) GetHasNewUserChat(req *types.GetHasNewUserChatReq) (resp *types.GetHasNewUserChatResp, err error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	toId := strconv.FormatInt(userId, 10)
	userChatListKey := globalkey.UserChatListKey + toId
	userChatList, err := l.svcCtx.RedisClient.ZrevrangeCtx(l.ctx, userChatListKey, 0, 0)
	if userChatList == nil || len(userChatList) == 0 {
		return &types.GetHasNewUserChatResp{
			HasNew: false,
		}, nil
	}
	getHasNewFriendChatLogic := NewGetHasNewFriendChatLogic(l.ctx, l.svcCtx)
	fromId, err := strconv.ParseInt(userChatList[0], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("fail to conv value"), "fail to conv value %v, err: %v", userChatList[0], err)
	}
	getHasNewFriendChatResp, err := getHasNewFriendChatLogic.GetHasNewFriendChat(&types.GetHasNewFriendChatReq{
		FromId: fromId,
		ToId:   userId,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetHasNewUserChatResp{
		HasNew: getHasNewFriendChatResp.HasNew,
	}, nil
}
