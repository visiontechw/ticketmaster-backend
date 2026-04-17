package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visiontechw/ticketmaster/internal/app/dto"
	usecases "github.com/visiontechw/ticketmaster/internal/app/usecases/user"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

type UserHandler struct {
	createUserUC  *usecases.CreateUserUseCase
	loginUserUC   *usecases.LoginUseCase
	getUserByIdUC *usecases.GetUserByIdUseCase
}

func NewUserHandler(uc *usecases.CreateUserUseCase,
	loginUC *usecases.LoginUseCase,
	getById *usecases.GetUserByIdUseCase) *UserHandler {
	return &UserHandler{
		createUserUC:  uc,
		loginUserUC:   loginUC,
		getUserByIdUC: getById,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.createUserUC.Execute(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrEmailAlreadyUsed {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (h *UserHandler) Me(c *gin.Context) {

	output, err := h.getUserByIdUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input dto.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.loginUserUC.Execute(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, output)
}
