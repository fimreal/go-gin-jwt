// Package classification Petstore API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//	Schemes: http, https
//	Host: localhost
//	BasePath: /account
//	Version: 0.0.1
//	License: MIT http://opensource.org/licenses/MIT
//	Contact: lmr@epurs.com
//
// swagger:meta
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
