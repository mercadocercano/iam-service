package port

import (
	"iam/src/shared/domain/criteria"
	"iam/src/role/domain/entity"
)

// RoleCriteriaRepository extiende RoleRepository con soporte para criteria
type RoleCriteriaRepository interface {
	RoleRepository
	criteria.CriteriaRepository[entity.Role]
}