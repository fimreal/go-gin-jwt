package rest

import (
	"github.com/fimreal/go-gin-jwt/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AddAccountRoutes(r *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = "/account"
	// account := r.Group("/account")
	account := r.Group(docs.SwaggerInfo.BasePath)
	h, _ := NewHandler()
	account.POST("/signin", h.AddUser)
	account.POST("/login", h.Login)
	account.POST("/valid", h.IsValid)
	account.POST("/renew", h.RenewToken)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
