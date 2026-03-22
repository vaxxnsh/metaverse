package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/handler"
)

type AppHandlers struct {
	handler.AuthHandler
}

func SetupRouter(appHandlers AppHandlers) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/user/signup", appHandlers.RegisterUser)
			auth.POST("/user/login", appHandlers.LoginUser)
			auth.POST("/admin/signup", appHandlers.RegisterAdmin)
			auth.POST("/admin/login", appHandlers.LoginAdmin)
		}
	}

	return r
}
