//go:build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Tipos de respuesta para Plans ---

type planResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Status      string   `json:"status"`
	MaxUsers    int      `json:"max_users"`
	PriceMonth  float64  `json:"price_month"`
	PriceYear   float64  `json:"price_year"`
	Features    []string `json:"features"`
}

type listPlansResponse struct {
	Items      []planResponse `json:"items"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// --- Tests de Plans ---

func TestPlans_POST_HappyPath_Returns201(t *testing.T) {
	srv := newTestServer(t)
	url := baseURL(srv) + "/plans"

	// El seed ya creó planes con tipos FREE/BASIC/PREMIUM/ENTERPRISE
	// Usamos un nombre único para evitar conflicto con el seed
	body := map[string]interface{}{
		"name":        "Plan Custom Integration",
		"description": "Plan de integración personalizado para testing",
		"type":        "BASIC",
		"price_month": 15.0,
		"price_year":  150.0,
		"features":    []string{"Custom Feature 1", "Custom Feature 2"},
	}
	resp := postJSON(t, url, body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var plan planResponse
	decodeJSON(t, resp, &plan)
	assert.Equal(t, "Plan Custom Integration", plan.Name)
	assert.Equal(t, "BASIC", plan.Type)
	assert.Equal(t, "ACTIVE", plan.Status)
	assert.NotEmpty(t, plan.ID)
	_, err := uuid.Parse(plan.ID)
	assert.NoError(t, err, "id debe ser UUID válido")
	assert.Equal(t, 15.0, plan.PriceMonth)
}

func TestPlans_POST_DuplicateName_Returns409(t *testing.T) {
	srv := newTestServer(t)
	url := baseURL(srv) + "/plans"

	uniquePlan := map[string]interface{}{
		"name":        "Plan Duplicado Test",
		"description": "Plan de duplicado para testing de integración",
		"type":        "ENTERPRISE",
		"price_month": 99.0,
		"price_year":  990.0,
	}
	postJSON(t, url, uniquePlan)
	resp := postJSON(t, url, uniquePlan)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestPlans_POST_InvalidType_Returns400(t *testing.T) {
	srv := newTestServer(t)
	url := baseURL(srv) + "/plans"

	body := map[string]interface{}{
		"name":        "Plan Tipo Invalido",
		"description": "Descripcion del plan de tipo invalido test",
		"type":        "SUPER_PLAN",
	}
	resp := postJSON(t, url, body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPlans_POST_MissingRequiredFields_Returns400(t *testing.T) {
	srv := newTestServer(t)
	url := baseURL(srv) + "/plans"

	resp := postJSON(t, url, map[string]interface{}{"name": "Solo nombre"})
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPlans_GET_ByID_HappyPath_Returns200(t *testing.T) {
	srv := newTestServer(t)
	base := baseURL(srv)

	// Obtener un plan existente del seed via listado
	listResp, err := http.Get(base + "/plans?page=1&page_size=10")
	require.NoError(t, err)
	defer listResp.Body.Close()

	var list listPlansResponse
	require.NoError(t, json.NewDecoder(listResp.Body).Decode(&list))
	require.Greater(t, len(list.Items), 0, "el seed debe haber creado planes")

	existingPlanID := list.Items[0].ID

	resp, err := http.Get(fmt.Sprintf("%s/plans/%s", base, existingPlanID))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched planResponse
	decodeJSON(t, resp, &fetched)
	assert.Equal(t, existingPlanID, fetched.ID)
	assert.NotEmpty(t, fetched.Name)
}

func TestPlans_GET_ByID_NotFound_Returns404(t *testing.T) {
	srv := newTestServer(t)
	url := fmt.Sprintf("%s/plans/%s", baseURL(srv), uuid.New().String())

	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestPlans_GET_ByID_InvalidUUID_Returns400(t *testing.T) {
	srv := newTestServer(t)
	url := baseURL(srv) + "/plans/not-a-uuid"

	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPlans_GET_List_ReturnsPaginationShape(t *testing.T) {
	srv := newTestServer(t)
	base := baseURL(srv)

	resp, err := http.Get(base + "/plans?page=1&page_size=10")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp listPlansResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&listResp))

	// El seed inserta 4 planes
	assert.GreaterOrEqual(t, listResp.TotalCount, 4)
	assert.Equal(t, 1, listResp.Page)
	assert.Equal(t, 10, listResp.PageSize)
	assert.Greater(t, listResp.TotalPages, 0)
	assert.GreaterOrEqual(t, len(listResp.Items), 4)
}

func TestPlans_GET_List_DefaultPagination_Returns200(t *testing.T) {
	srv := newTestServer(t)
	base := baseURL(srv)

	resp, err := http.Get(base + "/plans")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp listPlansResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&listResp))
	assert.GreaterOrEqual(t, listResp.TotalCount, 0)
	assert.Greater(t, listResp.PageSize, 0)
}

func TestPlans_GET_List_ContainsExpectedFields(t *testing.T) {
	srv := newTestServer(t)
	base := baseURL(srv)

	resp, err := http.Get(base + "/plans?page=1&page_size=10")
	require.NoError(t, err)
	defer resp.Body.Close()

	var listResp listPlansResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&listResp))
	require.Greater(t, len(listResp.Items), 0)

	// Verificar estructura del primer elemento
	first := listResp.Items[0]
	assert.NotEmpty(t, first.ID)
	assert.NotEmpty(t, first.Name)
	assert.NotEmpty(t, first.Type)
	assert.NotEmpty(t, first.Status)
	assert.GreaterOrEqual(t, first.MaxUsers, -1)
}
