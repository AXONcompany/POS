package user

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/domain/user"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

var roleNames = map[int]string{
	1: "ADMIN",
	2: "CAJA",
	3: "MESERO",
}

func toUsuarioResponse(u *user.User) gin.H {
	rol := roleNames[u.RoleID]
	if rol == "" {
		rol = "DESCONOCIDO"
	}

	resp := gin.H{
		"id":             u.ID,
		"nombre":         u.Name,
		"email":          u.Email,
		"rol":            rol,
		"activo":         u.IsActive,
		"fecha_creacion": u.CreatedAt,
	}
	if u.Phone != nil {
		resp["telefono"] = *u.Phone
	}
	if u.LastAccess != nil {
		resp["ultimo_acceso"] = *u.LastAccess
	}
	return resp
}

func (h *Handler) GetAll(c *gin.Context) {
	venueID, exists := c.Get(middleware.VenueIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("unauthorized", "UNAUTHORIZED"))
		return
	}

	users, err := h.uc.GetAllUsers(c.Request.Context(), venueID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	result := make([]gin.H, len(users))
	for i, u := range users {
		result[i] = toUsuarioResponse(u)
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(result))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	u, err := h.uc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Usuario no encontrado", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toUsuarioResponse(u)))
}

type UpdateUserRequest struct {
	Name   *string `json:"nombre"`
	Email  *string `json:"email"`
	RoleID *int    `json:"rol_id"`
	Active *bool   `json:"activo"`
	Phone  *string `json:"telefono"`
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	u, err := h.uc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Usuario no encontrado", "NOT_FOUND"))
		return
	}

	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.RoleID != nil {
		u.RoleID = *req.RoleID
	}
	if req.Active != nil {
		u.IsActive = *req.Active
	}
	if req.Phone != nil {
		u.Phone = req.Phone
	}

	updated, err := h.uc.UpdateUser(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toUsuarioResponse(updated)))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	if err := h.uc.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Usuario desactivado exitosamente"))
}
