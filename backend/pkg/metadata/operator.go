package metadata

import (
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

// NewSystemOperator 创建系统操作员的 OperatorMetadata
func NewSystemOperator(tenantID uint64) *authenticationV1.OperatorMetadata {
	return &authenticationV1.OperatorMetadata{
		Type:      authenticationV1.OperatorMetadata_SYSTEM,
		UserId:    0, // 系统操作员用户ID通常为0
		TenantId:  tenantID,
		DataScope: identityV1.DataScope_ALL, // 系统操作员拥有全部数据权限
	}
}

// NewServiceOperator 创建服务操作员的 OperatorMetadata
func NewServiceOperator(userID, tenantID uint64) *authenticationV1.OperatorMetadata {
	return &authenticationV1.OperatorMetadata{
		Type:      authenticationV1.OperatorMetadata_SERVICE,
		UserId:    userID,
		TenantId:  tenantID,
		DataScope: identityV1.DataScope_ALL, // 服务操作员拥有全部数据权限
	}
}

// NewUserOperator 创建用户操作员的 OperatorMetadata
func NewUserOperator(userID, tenantID, orgUnitID uint64, dataScope identityV1.DataScope) *authenticationV1.OperatorMetadata {
	return &authenticationV1.OperatorMetadata{
		Type:      authenticationV1.OperatorMetadata_USER,
		UserId:    userID,
		TenantId:  tenantID,
		OrgUnitId: orgUnitID,
		DataScope: dataScope,
	}
}
