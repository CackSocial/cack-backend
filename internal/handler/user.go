package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/usecase/user"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase *user.UserUseCase
}

func NewUserHandler(uc *user.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: uc}
}

func (h *UserHandler) RegisterRoutes(public, protected *gin.RouterGroup, optionalAuth gin.HandlerFunc) {
	public.POST("/auth/register", h.Register)
	public.POST("/auth/login", h.Login)
	public.GET("/users/:username", optionalAuth, h.GetProfile)

	protected.PUT("/users/me", h.UpdateProfile)
	protected.DELETE("/users/me", h.DeleteAccount)
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with username, password, and optional display name
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Register request"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.userUseCase.Register(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.userUseCase.Login(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get a user's public profile by username
// @Tags Users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{username} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	username := c.Param("username")
	currentUserID := getUserID(c)

	profile, err := h.userUseCase.GetProfile(username, currentUserID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Param body body dto.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserID(c)
	profile, err := h.userUseCase.UpdateProfile(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, profile)
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Permanently delete the authenticated user's account and all associated data
// @Tags Users
// @Accept json
// @Produce json
// @Param body body dto.DeleteAccountRequest true "Password confirmation"
// @Success 204
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /users/me [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	var req dto.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserID(c)
	if err := h.userUseCase.DeleteAccount(userID, req.Password); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
