package metadata

import (
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

// NewUserOperator 创建用户操作员的 OperatorMetadata
func NewUserOperator(userID, tenantID, orgUnitID int64, dataScope identityV1.DataScope) *authenticationV1.OperatorMetadata {
	return &authenticationV1.OperatorMetadata{
		Type:      authenticationV1.OperatorMetadata_USER,
		UserId:    userID,
		TenantId:  tenantID,
		OrgUnitId: orgUnitID,
		DataScope: dataScope,
	}
}
