package app

import (
	"github.com/fimreal/go-gin-jwt/rest"
	"github.com/gin-gonic/gin"
)

func Run(address string) error {
	r := gin.Default()
	apiv1 := r.Group("/")
	rest.AddAccountRoutes(apiv1)
	return r.Run(address)
}
