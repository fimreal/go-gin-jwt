package rest

import "github.com/gin-gonic/gin"

func AddAccountRoutes(r *gin.RouterGroup) {
	account := r.Group("/account")
	h, _ := NewHandler()
	account.POST("/signin", h.AddUser)
	account.POST("/login", h.Login)
	account.POST("/valid", h.IsValid)
	account.POST("/renew", h.RenewToken)
}
