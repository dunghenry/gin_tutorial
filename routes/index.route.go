package routes

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func Run(r *gin.Engine) {
	port := os.Getenv("PORT")
	// r.LoadHTMLGlob("templates/*")
	r.LoadHTMLGlob("templates/**/*")
	getRoutes(r)
	fmt.Printf("Server listen on:http://localhost:")
	fmt.Printf(port)
	r.Run(":3000")
}

func getRoutes(r *gin.Engine) {
	route := r.Group("/")
	api := r.Group("/api")
	SiteRoutes(route)
	AuthRoutes(route)
	UserRoutes(api)
}
