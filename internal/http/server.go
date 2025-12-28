package http

import (
	"github.com/AXONcompany/POS/internal/config"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	RegisterRouters(r)
	return r

}
