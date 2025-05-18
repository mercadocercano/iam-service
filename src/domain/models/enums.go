package models

type SaasType string
type UserStatus string

const (
	CRM       SaasType = "CRM"
	ERP       SaasType = "ERP"
	ECOMMERCE SaasType = "ECOMMERCE"
	ALL       SaasType = "ALL"
)

func (s SaasType) IsValid() bool {
	switch s {
	case CRM, ERP, ECOMMERCE, ALL:
		return true
	default:
		return false
	}
}

const (
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
	StatusPending  UserStatus = "PENDING"
)
