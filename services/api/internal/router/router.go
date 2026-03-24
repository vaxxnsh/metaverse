package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/handler"
	"github.com/vaxxnsh/metaverse/api/internal/middleware"
)

type AppHandlers struct {
	handler.AuthHandler
	handler.UserHandler
	handler.AdminHandler
	handler.MetadataHandler
	handler.SpaceHandler
}

func SetupRouter(appHandlers AppHandlers) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/user/signup", appHandlers.RegisterUser)
				auth.POST("/user/login", appHandlers.LoginUser)
				auth.POST("/admin/signup", appHandlers.RegisterAdmin)
				auth.POST("/admin/login", appHandlers.LoginAdmin)
			}

			v1.PATCH("/user/metadata", middleware.AuthUser(), appHandlers.UserHandler.PatchMetadata)
			v1.PATCH("/admin/metadata", middleware.AuthAdmin(), appHandlers.AdminHandler.PatchMetadata)
			v1.GET("/user/metadata/bulk", middleware.AuthUser(), appHandlers.MetadataHandler.GetBulkUserAvatars)

			v1.POST("/space", middleware.AuthUser(), appHandlers.SpaceHandler.CreateSpace)
		}
	}

	return r
}
