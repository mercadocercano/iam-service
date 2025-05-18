package handlers

import (
	"log"
	"iam/src/application/services"
	"iam/src/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PlanHandler struct {
	planService *services.PlanService
}

func NewPlanHandler(planService *services.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
	}
}

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.planService.CreatePlan(c.Request.Context(), &plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *PlanHandler) GetPlanByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	plan, err := h.planService.GetPlanByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *PlanHandler) GetAllPlans(c *gin.Context) {
	saasType := c.Query("saas")
	var plans []models.Plan
	var err error

	if saasType != "" {
		plans, err = h.planService.GetPlansBySaas(c.Request.Context(), models.SaasType(saasType))
	} else {
		plans, err = h.planService.GetAllPlans(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Primero obtener el plan existente
	existingPlan, err := h.planService.GetPlanByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Crear una copia del plan existente
	updatedPlan := *existingPlan

	// Estructura temporal para recibir solo los campos actualizables
	var updateData struct {
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		Features     []string `json:"features"`
		MonthlyPrice float64  `json:"monthly_price"`
		YearlyPrice  float64  `json:"yearly_price"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar solo los campos permitidos
	updatedPlan.Name = updateData.Name
	updatedPlan.Description = updateData.Description
	features := pq.StringArray(updateData.Features)
	updatedPlan.Features = &features
	updatedPlan.MonthlyPrice = updateData.MonthlyPrice
	updatedPlan.YearlyPrice = updateData.YearlyPrice

	log.Printf("Plan before update - ID: %v, Saas: %v, Name: %v", updatedPlan.ID, updatedPlan.Saas, updatedPlan.Name)

	if err := h.planService.UpdatePlan(c.Request.Context(), &updatedPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPlan)
}

func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.planService.DeletePlan(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
