package httputil

import "github.com/gin-gonic/gin"

// SuccessResponse devuelve el wrapper estandarizado para respuestas exitosas.
func SuccessResponse(data interface{}) gin.H {
	return gin.H{
		"success": true,
		"data":    data,
	}
}

// SuccessMessageResponse devuelve una respuesta exitosa con mensaje.
func SuccessMessageResponse(message string) gin.H {
	return gin.H{
		"success": true,
		"message": message,
	}
}

// ErrorResponse devuelve el wrapper estandarizado para errores.
func ErrorResponse(message, code string) gin.H {
	resp := gin.H{
		"success": false,
		"error":   message,
	}
	if code != "" {
		resp["code"] = code
	}
	return resp
}
