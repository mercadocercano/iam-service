package port

import (
	"iam/src/shared/domain/criteria"
	"iam/src/plan/domain/entity"
)

// PlanCriteriaRepository extiende PlanRepository con soporte para criteria
type PlanCriteriaRepository interface {
	PlanRepository
	criteria.CriteriaRepository[entity.Plan]
}