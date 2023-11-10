package routes

import (
	"gin-wire-template/controllers"
	"github.com/gin-gonic/gin"
)

func AddApi(r *gin.Engine, cs *controllers.Controllers) *gin.Engine {
	api := r.Group("/api/v1")
	{
		api.GET("/text_menu", cs.TextMenuController.GetTextMenu)
		api.POST("/text_menu", cs.TextMenuController.CreateTextMenu)
		api.PUT("/text_menu/:id", cs.TextMenuController.UpdateTextMenu)
		api.DELETE("/text_menu/:id", cs.TextMenuController.DeleteTextMenu)
	}

	return r
}
