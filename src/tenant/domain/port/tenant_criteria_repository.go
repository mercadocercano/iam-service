package port

import (
	"iam/src/shared/domain/criteria"
	"iam/src/tenant/domain/entity"
)

// TenantCriteriaRepository extiende TenantRepository con soporte para criteria
type TenantCriteriaRepository interface {
	TenantRepository
	criteria.CriteriaRepository[entity.Tenant]
}