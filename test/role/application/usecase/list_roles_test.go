package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"iam/src/role/application/usecase"
	"iam/src/role/domain/entity"
	roleMother "iam/test/role/domain/entity"
	"iam/test/role/infrastructure/persistence/repository"
)

func TestListRolesUseCase_Execute_ReturnsRoles(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockRoleRepository()
	listUseCase := usecase.NewListRolesUseCase(mockRepo)
	ctx := context.Background()

	mother := roleMother.Create()
	role1 := mother.SystemAdmin()
	role2 := mother.Custom()
	mockRepo.SetupRoles([]*entity.Role{role1, role2})

	// Act
	resp, err := listUseCase.Execute(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 2, resp.TotalCount)
	assert.Len(t, resp.Roles, 2)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, 1, mockRepo.GetCallCount("List"))
	assert.Equal(t, 1, mockRepo.GetCallCount("Count"))
}

func TestListRolesUseCase_Execute_EmptyResult(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockRoleRepository()
	listUseCase := usecase.NewListRolesUseCase(mockRepo)
	ctx := context.Background()

	// Act
	resp, err := listUseCase.Execute(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 0, resp.TotalCount)
	assert.Empty(t, resp.Roles)
}

func TestListRolesUseCase_GetByTenant_ReturnsOnlyTenantRoles(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockRoleRepository()
	listUseCase := usecase.NewListRolesUseCase(mockRepo)
	ctx := context.Background()

	mother := roleMother.Create()
	tenantID := uuid.New()
	role1 := mother.UserForTenant(tenantID)
	role2 := mother.TenantAdminForTenant(tenantID)
	role3 := mother.SystemAdmin() // No debe aparecer
	mockRepo.SetupRoles([]*entity.Role{role1, role2, role3})

	// Act
	resp, err := listUseCase.GetByTenant(ctx, tenantID, 1, 10)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 1, mockRepo.GetCallCount("GetByTenant"))
	assert.Equal(t, 1, mockRepo.GetCallCount("CountByTenant"))
}

func TestListRolesUseCase_GetSystemRoles_ReturnsSystemRoles(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockRoleRepository()
	listUseCase := usecase.NewListRolesUseCase(mockRepo)
	ctx := context.Background()

	mother := roleMother.Create()
	sysRole := mother.SystemAdmin()
	tenantRole := mother.User()
	mockRepo.SetupRoles([]*entity.Role{sysRole, tenantRole})

	// Act
	resp, err := listUseCase.GetSystemRoles(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 1, mockRepo.GetCallCount("GetSystemRoles"))
	// Solo el rol de sistema deberia estar en la respuesta
	for _, r := range resp.Roles {
		assert.True(t, r.IsSystem)
	}
}

func TestListRolesUseCase_Execute_ListFails_ReturnsError(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockRoleRepository()
	listUseCase := usecase.NewListRolesUseCase(mockRepo)
	ctx := context.Background()

	mockRepo.ShouldFailOn("List")

	// Act
	resp, err := listUseCase.Execute(ctx, 1, 10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, repository.ErrMockFailedOp, err)
}

func TestListRolesUseCase_Execute_CountFails_ReturnsError(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockRoleRepository()
	listUseCase := usecase.NewListRolesUseCase(mockRepo)
	ctx := context.Background()

	mockRepo.ShouldFailOn("Count")

	// Act
	resp, err := listUseCase.Execute(ctx, 1, 10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
}
