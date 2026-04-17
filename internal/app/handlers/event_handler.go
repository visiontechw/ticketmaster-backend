package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/app/dto"
	usecases "github.com/visiontechw/ticketmaster/internal/app/usecases/events"
)

type EventHandler struct {
	getEventsUC   *usecases.GetEventsUseCase
	getEventById  *usecases.GetEventByIdUseCase
	createEventUC *usecases.CreateventUseCase
	updateEventUC *usecases.UpdateEventUseCase
	deleteEventUC *usecases.DeleteEventUseCase
}

func NewEventHandler(geuc *usecases.GetEventsUseCase, ceuc *usecases.CreateventUseCase,
	geid *usecases.GetEventByIdUseCase, ueuc *usecases.UpdateEventUseCase, deuc *usecases.DeleteEventUseCase) *EventHandler {
	return &EventHandler{
		getEventsUC:   geuc,
		createEventUC: ceuc,
		getEventById:  geid,
		updateEventUC: ueuc,
		deleteEventUC: deuc,
	}
}

func (h *EventHandler) Create(c *gin.Context) {
	var input dto.CreateEventRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.createEventUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, event)

}

func (h *EventHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	eventId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var input dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.updateEventUC.Execute(c.Request.Context(), eventId, input)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	eventId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	// Extract user_id from context
	userId, ok := c.Request.Context().Value("user_id").(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: user not found"})
		return
	}

	err = h.deleteEventUC.Execute(c.Request.Context(), eventId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") || strings.Contains(err.Error(), "não tem permissão") {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found or you don't have permission to delete"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}

func (h *EventHandler) GetById(c *gin.Context) {
	idParam := c.Param("id")

	eventId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de ID inválido"})
		return
	}

	event, err := h.getEventById.Execute(c.Request.Context(), eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "evento não encontrado"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	// 1. Extrair e validar Query Parameters
	activeParam := c.DefaultQuery("active", "true")
	isActive := activeParam != "false"

	// 2. Montar o Input para o Usecase
	input := dto.ListEventsInput{
		Active: isActive,
	}

	events, err := h.getEventsUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch events",
		})
		return
	}

	c.JSON(http.StatusOK, events)
}
