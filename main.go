package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/mbarolo/webp_backend/config"
	_ "github.com/mbarolo/webp_backend/docs"
	"github.com/mbarolo/webp_backend/routes"
	"github.com/mbarolo/webp_backend/utils"
)

func main() {

	// se cargan las variables de entorno
	utils.LoadEnv()

	// logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Start test_back")
	log.Printf("serverUp, %s", os.Getenv("ADDR"))

	defer config.CloseDB()

	app := gin.Default()
	app.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	)))

	// 404 por defecto
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Servicio no encontrado."})
	})

	routes.InitRoutes(app)

	for _, route := range app.Routes() {
		log.Printf("[%s] %s", route.Method, route.Path)
	}

	log.Printf("Server starting on %s", os.Getenv("ADDR"))
	app.Run(os.Getenv("ADDR"))
}
