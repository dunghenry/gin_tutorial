package routes

import (
	"trandung/gin_tutorial/controllers"

	"github.com/gin-gonic/gin"
)

func SiteRoutes(rg *gin.RouterGroup) {
	site := rg.Group("/")
	site.GET("/", controllers.GetHomePage)
	site.GET("/upload", controllers.GetUploadPage)
	site.GET("/posts", controllers.GetPostPage)
	site.POST("/upload", controllers.PostUpload)
}
