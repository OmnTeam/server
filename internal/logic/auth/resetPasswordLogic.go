package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/perfect-panel/server/internal/model/log"
	"github.com/perfect-panel/server/internal/model/user"
	"github.com/perfect-panel/server/pkg/jwt"
	"github.com/perfect-panel/server/pkg/uuidx"

	"github.com/perfect-panel/server/internal/config"
	"github.com/perfect-panel/server/pkg/constant"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/tool"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type ResetPasswordLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewResetPasswordLogic Reset password
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordRequest) (resp *types.LoginResponse, err error) {
	var userInfo *user.User
	loginStatus := false

	defer func() {
		if userInfo.Id != 0 && loginStatus {
			loginLog := log.Login{
				Method:    "email",
				LoginIP:   req.IP,
				UserAgent: req.UserAgent,
				Success:   loginStatus,
				Timestamp: time.Now().UnixMilli(),
			}
			content, _ := loginLog.Marshal()
			if err := l.svcCtx.LogModel.Insert(l.ctx, &log.SystemLog{
				Id:       0,
				Type:     log.TypeLogin.Uint8(),
				Date:     time.Now().Format("2006-01-02"),
				ObjectID: userInfo.Id,
				Content:  string(content),
			}); err != nil {
				l.Errorw("failed to insert login log",
					logger.Field("user_id", userInfo.Id),
					logger.Field("ip", req.IP),
					logger.Field("error", err.Error()),
				)
			}
		}
	}()

	cacheKey := fmt.Sprintf("%s:%s:%s", config.AuthCodeCacheKey, constant.Security, req.Email)
	// Check the verification code
	if value, err := l.svcCtx.Redis.Get(l.ctx, cacheKey).Result(); err != nil {
		l.Errorw("Verification code error", logger.Field("cacheKey", cacheKey), logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "Verification code error")
	} else {
		var payload CacheKeyPayload
		if err := json.Unmarshal([]byte(value), &payload); err != nil {
			l.Errorw("Unmarshal errors", logger.Field("cacheKey", cacheKey), logger.Field("error", err.Error()), logger.Field("value", value))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "Verification code error")
		}
		if payload.Code != req.Code {
			l.Errorw("Verification code error", logger.Field("cacheKey", cacheKey), logger.Field("error", "Verification code error"), logger.Field("reqCode", req.Code), logger.Field("payloadCode", payload.Code))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.VerifyCodeError), "Verification code error")
		}
	}

	// Check user
	authMethod, err := l.svcCtx.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user email not exist: %v", req.Email)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find user by email error: %v", err.Error())
	}

	userInfo, err = l.svcCtx.UserModel.FindOne(l.ctx, authMethod.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user email not exist: %v", req.Email)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query user info failed: %v", err.Error())
	}

	// Update password
	userInfo.Password = tool.EncodePassWord(req.Password)
	if err := l.svcCtx.UserModel.Update(l.ctx, userInfo); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update user info failed: %v", err.Error())
	}

	loginType := "pc"

	// Bind device to user if identifier is provided
	if req.Identifier != "" {
		bindLogic := NewBindDeviceLogic(l.ctx, l.svcCtx)
		if err := bindLogic.BindDeviceToUser(req.Identifier, req.IP, req.UserAgent, userInfo.Id); err != nil {
			l.Errorw("failed to bind device to user",
				logger.Field("user_id", userInfo.Id),
				logger.Field("identifier", req.Identifier),
				logger.Field("error", err.Error()),
			)
			// Don't fail register if device binding fails, just log the error
		}
		loginType = "mobile"
	}

	// Generate session id
	sessionId := uuidx.NewUUID().String()
	// Generate token
	token, err := jwt.NewJwtToken(
		l.svcCtx.Config.JwtAuth.AccessSecret,
		time.Now().Unix(),
		l.svcCtx.Config.JwtAuth.AccessExpire,
		jwt.WithOption("UserId", userInfo.Id),
		jwt.WithOption("SessionId", sessionId),
		jwt.WithOption("LoginType", loginType),
	)
	if err != nil {
		l.Logger.Error("[UserLogin] token generate error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "token generate error: %v", err.Error())
	}
	sessionIdCacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
	if err = l.svcCtx.Redis.Set(l.ctx, sessionIdCacheKey, userInfo.Id, time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire)*time.Second).Err(); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "set session id error: %v", err.Error())
	}
	loginStatus = true
	return &types.LoginResponse{
		Token: token,
	}, nil
}
