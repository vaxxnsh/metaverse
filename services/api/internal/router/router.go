package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/user"
)

type AppHandlers struct {
	UserHandler user.Handler
}

func SetupRouter(appHandlers AppHandlers) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/signup", appHandlers.UserHandler.Signup)
			users.POST("/", appHandlers.UserHandler.Login)
		}
	}

	return r
}
