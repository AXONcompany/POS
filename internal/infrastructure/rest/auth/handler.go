package auth

import (
	"net/http"
	"time"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	"github.com/AXONcompany/POS/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

var roleNames = map[int]string{
	1: "PROPIETARIO",
	2: "CAJERO",
	3: "MESERO",
}

type Handler struct {
	uc *auth.Usecase
}

func NewHandler(uc *auth.Usecase) *Handler {
	return &Handler{uc: uc}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := c.GetHeader("User-Agent")

	tokens, err := h.uc.Login(c.Request.Context(), req.Email, req.Password, deviceInfo, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("Email o contrasena incorrectos", "UNAUTHORIZED"))
		return
	}

	u := tokens.User
	c.SetCookie("refresh_token", tokens.RefreshToken, int(24*time.Hour.Seconds()), "/", "", true, true)

	rol := roleNames[u.RoleID]
	if rol == "" {
		rol = "DESCONOCIDO"
	}

	userData := gin.H{
		"id":             u.ID,
		"nombre":         u.Name,
		"email":          u.Email,
		"rol":            rol,
		"activo":         u.IsActive,
		"fecha_creacion": u.CreatedAt,
	}
	if u.Phone != nil {
		userData["telefono"] = *u.Phone
	}
	if u.LastAccess != nil {
		userData["ultimo_acceso"] = *u.LastAccess
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"token":      tokens.AccessToken,
		"usuario":    userData,
		"expires_in": 900,
	}))
}

type RegisterRequest struct {
	Name     string `json:"nombre" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"rol" binding:"required"`
	Phone    string `json:"telefono"`
}

// Register crea un usuario. Propietarios pueden crear cualquier rol.
// Cajeros solo pueden crear meseros.
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	roleID := 0
	switch req.Role {
	case "PROPIETARIO":
		roleID = 1
	case "CAJERO":
		roleID = 2
	case "MESERO":
		roleID = 3
	default:
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Rol invalido, debe ser PROPIETARIO, CAJERO o MESERO", "BAD_REQUEST"))
		return
	}

	// Verificar permisos: cajeros solo pueden crear meseros
	callerRoleID, _ := c.Get(middleware.RoleIDKey)
	if callerRoleID.(int) == 2 && roleID != 3 { // Cajero solo puede crear meseros
		c.JSON(http.StatusForbidden, httputil.ErrorResponse("Cajeros solo pueden registrar meseros", "FORBIDDEN"))
		return
	}

	venueID, exists := c.Get(middleware.VenueIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("No autorizado", "UNAUTHORIZED"))
		return
	}

	created, err := h.uc.RegisterUser(c.Request.Context(), req.Name, req.Email, req.Password, roleID, venueID.(int), req.Phone)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, httputil.ErrorResponse("El email ya esta registrado en el sistema", "CONFLICT"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	rol := roleNames[created.RoleID]

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Usuario registrado exitosamente",
		"data": gin.H{
			"id":             created.ID,
			"nombre":         created.Name,
			"email":          created.Email,
			"rol":            rol,
			"activo":         created.IsActive,
			"fecha_creacion": created.CreatedAt,
		},
	})
}

// RegisterOwner registra un nuevo propietario con su primera sede (endpoint publico).
type RegisterOwnerRequest struct {
	Name      string `json:"nombre" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	VenueName string `json:"nombre_sede" binding:"required"`
	Address   string `json:"direccion"`
	Phone     string `json:"telefono"`
}

func (h *Handler) RegisterOwner(c *gin.Context) {
	var req RegisterOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := c.GetHeader("User-Agent")

	tokens, err := h.uc.RegisterOwnerWithVenue(
		c.Request.Context(),
		req.Name, req.Email, req.Password,
		req.VenueName, req.Address, req.Phone,
		deviceInfo, ipAddress,
	)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, httputil.ErrorResponse("El email ya esta registrado", "CONFLICT"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	u := tokens.User
	c.SetCookie("refresh_token", tokens.RefreshToken, int(24*time.Hour.Seconds()), "/", "", true, true)

	c.JSON(http.StatusCreated, httputil.SuccessResponse(gin.H{
		"token": tokens.AccessToken,
		"usuario": gin.H{
			"id":     u.ID,
			"nombre": u.Name,
			"email":  u.Email,
			"rol":    "PROPIETARIO",
		},
		"message":    "Propietario y sede creados exitosamente",
		"expires_in": 900,
	}))
}

// RegisterWaiter crea un mesero con credenciales auto-generadas.
type RegisterWaiterRequest struct {
	Name  string `json:"nombre" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) RegisterWaiter(c *gin.Context) {
	var req RegisterWaiterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	venueID, exists := c.Get(middleware.VenueIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("No autorizado", "UNAUTHORIZED"))
		return
	}

	created, rawPassword, err := h.uc.RegisterWaiter(c.Request.Context(), req.Name, req.Email, venueID.(int))
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, httputil.ErrorResponse("El email ya esta registrado", "CONFLICT"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, httputil.SuccessResponse(gin.H{
		"message": "Mesero creado exitosamente. Guarde las credenciales, no se mostraran de nuevo.",
		"credenciales": gin.H{
			"email":    created.Email,
			"password": rawPassword,
		},
		"usuario": gin.H{
			"id":     created.ID,
			"nombre": created.Name,
			"email":  created.Email,
			"rol":    "MESERO",
			"activo": created.IsActive,
		},
	}))
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("Token invalido", "UNAUTHORIZED"))
		return
	}

	u, err := h.uc.GetUserByID(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Usuario no encontrado", "NOT_FOUND"))
		return
	}

	rol := roleNames[u.RoleID]
	if rol == "" {
		rol = "DESCONOCIDO"
	}

	userData := gin.H{
		"id":             u.ID,
		"nombre":         u.Name,
		"email":          u.Email,
		"rol":            rol,
		"activo":         u.IsActive,
		"fecha_creacion": u.CreatedAt,
	}
	if u.Phone != nil {
		userData["telefono"] = *u.Phone
	}
	if u.LastAccess != nil {
		userData["ultimo_acceso"] = *u.LastAccess
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(userData))
}

func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Sesion cerrada exitosamente"))
		return
	}

	_ = h.uc.RevokeSession(c.Request.Context(), refreshToken)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Sesion cerrada exitosamente"))
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh renueva el access token usando el refresh token.
// El refresh token se puede enviar como cookie o en el body JSON.
func (h *Handler) Refresh(c *gin.Context) {
	// Intentar obtener de cookie primero
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		// Fallback al body
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Refresh token requerido", "BAD_REQUEST"))
		return
	}

	tokens, err := h.uc.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		// Limpiar cookie invalida
		c.SetCookie("refresh_token", "", -1, "/", "", true, true)
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("Token invalido o expirado", "UNAUTHORIZED"))
		return
	}

	u := tokens.User
	// Setear nueva cookie con el nuevo refresh token
	c.SetCookie("refresh_token", tokens.RefreshToken, int(24*time.Hour.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"token":      tokens.AccessToken,
		"expires_in": 900,
		"usuario": gin.H{
			"id":     u.ID,
			"nombre": u.Name,
			"email":  u.Email,
			"rol":    roleNames[u.RoleID],
		},
	}))
}

type SwitchSedeRequest struct {
	SedeID int `json:"sede_id" binding:"required"`
}

// SwitchSede maneja POST /auth/switch-sede.
func (h *Handler) SwitchSede(c *gin.Context) {
	var req SwitchSedeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Dato inválido: sede_id es requerido", "BAD_REQUEST"))
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("unauthorized", "UNAUTHORIZED"))
		return
	}

	deviceInfo := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	resp, err := h.uc.SwitchSede(c.Request.Context(), userID.(int), req.SedeID, deviceInfo, ipAddress)
	if err != nil {
		c.JSON(http.StatusForbidden, httputil.ErrorResponse(err.Error(), "FORBIDDEN"))
		return
	}

	secure := false
	if c.Request.TLS != nil {
		secure = true
	}
	c.SetCookie("refresh_token", resp.RefreshToken, int(7*24*time.Hour.Seconds()), "/", "", secure, true)

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"token":      resp.AccessToken,
		"expires_in": 900,
	}))
}
