package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	jwtutil "github.com/tx7do/go-utils/jwtutil"
	"github.com/tx7do/go-utils/trans"

	authn "github.com/tx7do/kratos-authn/engine"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

const (
	ClaimFieldUserName      = authn.ClaimFieldSubject // 用户名
	ClaimFieldUserID        = "uid"                   // 用户 ID
	ClaimFieldTenantID      = "tid"                   // 租户 ID
	ClaimFieldClientID      = "cid"                   // 客户端 ID
	ClaimFieldDeviceID      = "did"                   // 设备 ID
	ClaimFieldRoleCodes     = "roc"                   // 角色码列表
	ClaimFieldDataScope     = "ds"                    // 数据范围
	ClaimFieldOrgUnitID     = "ouid"                  // 组织单元 ID
	ClaimFieldPlatformAdmin = "ipa"                   // 平台超级管理员标记（与 proto json_name 同名）
)

const (
	// DefaultTokenExpiration 默认令牌过期时间：2 小时
	DefaultTokenExpiration = 2 * time.Hour

	// DefaultRefreshTokenExpiration 默认刷新令牌过期时间：7 天
	DefaultRefreshTokenExpiration = 7 * 24 * time.Hour
)

// NewUserTokenPayload 创建用户令牌
func NewUserTokenPayload(
	username string,
	userID int64,
	tenantID int64,
	orgUnitID *int64,
	roleCodes []string,
	dataScope *identityV1.DataScope,
	clientID *string,
	deviceID *string,
) *authenticationV1.UserTokenPayload {
	return &authenticationV1.UserTokenPayload{
		Username:  trans.Ptr(username),
		UserId:    userID,
		TenantId:  trans.Ptr(tenantID),
		OrgUnitId: orgUnitID,
		Roles:     roleCodes,
		ClientId:  clientID,
		DeviceId:  deviceID,
		DataScope: dataScope,
	}
}

// NewUserTokenAuthClaims 创建用户令牌认证声明
func NewUserTokenAuthClaims(
	tokenPayload *authenticationV1.UserTokenPayload,
	expirationTime *time.Time,
) *authn.AuthClaims {
	authClaims := authn.AuthClaims{
		ClaimFieldUserName:       tokenPayload.GetUsername(),
		ClaimFieldUserID:         tokenPayload.GetUserId(),
		ClaimFieldTenantID:       tokenPayload.GetTenantId(),
		authn.ClaimFieldIssuedAt: time.Now().Unix(),
	}

	if expirationTime != nil {
		authClaims[authn.ClaimFieldExpirationTime] = expirationTime.Unix()
	}

	if tokenPayload.Jti != nil {
		authClaims[authn.ClaimFieldJwtID] = tokenPayload.GetJti()
	}

	if len(tokenPayload.Roles) > 0 {
		authClaims[ClaimFieldRoleCodes] = tokenPayload.Roles
	}
	if tokenPayload.DeviceId != nil {
		authClaims[ClaimFieldDeviceID] = tokenPayload.GetDeviceId()
	}
	if tokenPayload.ClientId != nil {
		authClaims[ClaimFieldClientID] = tokenPayload.GetClientId()
	}

	if tokenPayload.DataScope != nil {
		authClaims[ClaimFieldDataScope] = tokenPayload.GetDataScope().String()
	}
	if tokenPayload.OrgUnitId != nil {
		authClaims[ClaimFieldOrgUnitID] = tokenPayload.GetOrgUnitId()
	}

	// 平台管理员标记以字符串形式入 claim（AuthClaims 值经 JSON 编解码，布尔直存
	// 会被某些中间件还原成 float64，字符串口径两端都稳定）
	if tokenPayload.IsPlatformAdmin != nil {
		authClaims[ClaimFieldPlatformAdmin] = strconv.FormatBool(tokenPayload.GetIsPlatformAdmin())
	}

	return &authClaims
}

// NewUserTokenPayloadWithClaims 从认证声明创建用户令牌
func NewUserTokenPayloadWithClaims(claims *authn.AuthClaims) (*authenticationV1.UserTokenPayload, error) {
	payload := &authenticationV1.UserTokenPayload{}

	sub, err := claims.GetSubject()
	if err != nil {
		log.Errorf("GetSubject failed: %v", err)
	}
	if sub != "" {
		payload.Username = trans.Ptr(sub)
	}

	jti, err := claims.GetJwtID()
	if err != nil {
		log.Errorf("GetJwtID failed: %v", err)
	}
	if jti != "" {
		payload.Jti = trans.Ptr(jti)
	}

	userId, err := claims.GetInt64(ClaimFieldUserID)
	if err != nil {
		log.Errorf("GetInt64 ClaimFieldUserID failed: %v", err)
	}
	if userId != 0 {
		payload.UserId = userId
	}

	tenantId, err := claims.GetInt64(ClaimFieldTenantID)
	if err != nil {
		log.Errorf("GetInt64 ClaimFieldTenantID failed: %v", err)
	}
	if tenantId != 0 {
		payload.TenantId = trans.Ptr(tenantId)
	}

	clientId, err := claims.GetString(ClaimFieldClientID)
	if err != nil {
		log.Errorf("GetString ClaimFieldClientID failed: %v", err)
	}
	if clientId != "" {
		payload.ClientId = trans.Ptr(clientId)
	}

	deviceId, err := claims.GetString(ClaimFieldDeviceID)
	if err != nil {
		log.Errorf("GetString ClaimFieldDeviceID failed: %v", err)
	}
	if deviceId != "" {
		payload.DeviceId = trans.Ptr(deviceId)
	}

	roleCodes, err := claims.GetStrings(ClaimFieldRoleCodes)
	if err != nil {
		log.Errorf("GetStrings ClaimFieldRoleCodes failed: %v", err)
	}
	if roleCodes != nil {
		payload.Roles = roleCodes
	}

	dataScope, err := claims.GetString(ClaimFieldDataScope)
	if err != nil {
		log.Errorf("GetString ClaimFieldDataScope failed: %v", err)
	}
	if dataScope != "" {
		v, ok := identityV1.DataScope_value[dataScope]
		if ok {
			payload.DataScope = trans.Ptr(identityV1.DataScope(v))
		}
	}

	orgUnitID, err := claims.GetInt64(ClaimFieldOrgUnitID)
	if err != nil {
		log.Errorf("GetInt64 ClaimFieldOrgUnitID failed [%v]", err)
	}
	if orgUnitID != 0 {
		payload.OrgUnitId = trans.Ptr(orgUnitID)
	}

	// 平台管理员标记还原：缺省（旧令牌）保持 nil=false
	if platformAdmin, err := claims.GetString(ClaimFieldPlatformAdmin); err == nil && platformAdmin != "" {
		if v, parseErr := strconv.ParseBool(platformAdmin); parseErr == nil {
			payload.IsPlatformAdmin = trans.Ptr(v)
		}
	}

	return payload, nil
}

// NewUserTokenPayloadWithJwtMapClaims 从 JWT MapClaims 创建用户令牌
func NewUserTokenPayloadWithJwtMapClaims(claims jwt.MapClaims) (*authenticationV1.UserTokenPayload, error) {
	payload := &authenticationV1.UserTokenPayload{}

	sub, err := claims.GetSubject()
	if err != nil {
		log.Errorf("GetSubject failed: %v", err)
	}
	if sub != "" {
		payload.Username = trans.Ptr(sub)
	}

	// 注意：本函数会被审计日志中间件的 extractAuthToken 调用，而那里仅解析 token
	// 不做签名校验。攻击者可构造签名错误但结构合法的恶意 token，若 claims 值类型
	// 不符（如把数字字段填成字符串），裸类型断言会 panic，导致中间件崩溃。
	// 因此所有断言必须用 ,ok 形式，类型不符时跳过该字段而不是 panic。
	if userId, ok := claims[ClaimFieldUserID].(float64); ok {
		payload.UserId = int64(userId)
	}

	if tenantId, ok := claims[ClaimFieldTenantID].(float64); ok {
		payload.TenantId = trans.Ptr(int64(tenantId))
	}

	if clientId, ok := claims[ClaimFieldClientID].(string); ok {
		payload.ClientId = trans.Ptr(clientId)
	}

	if deviceId, ok := claims[ClaimFieldDeviceID].(string); ok {
		payload.DeviceId = trans.Ptr(deviceId)
	}

	if dataScope, ok := claims[ClaimFieldDataScope].(string); ok {
		if v, ok2 := identityV1.DataScope_value[dataScope]; ok2 {
			payload.DataScope = trans.Ptr(identityV1.DataScope(v))
		}
	}

	if orgUnitID, ok := claims[ClaimFieldOrgUnitID].(float64); ok {
		payload.OrgUnitId = trans.Ptr(int64(orgUnitID))
	}

	roleCodes, _ := claims[ClaimFieldRoleCodes]
	if roleCodes != nil {
		switch itf := roleCodes.(type) {
		case []interface{}:
			for _, rc := range itf {
				if s, ok := rc.(string); ok {
					payload.Roles = append(payload.Roles, s)
				}
			}

		case []string:
			payload.Roles = itf

		default:
			return nil, errors.New("invalid roleCodes type")
		}
	}

	return payload, nil
}

// IsTokenExpired 检查令牌是否过期（判定口径委托库 jwtutil：nil 即过期、
// 缺 exp 不视为过期、60s leeway；库 claims_expiry.go 注释明写与本项目
// 历史实现约定一致，勿单独漂移）。
func IsTokenExpired(claims *authn.AuthClaims) bool {
	if claims == nil {
		return true
	}
	return jwtutil.IsTokenExpired(jwt.MapClaims(*claims))
}
