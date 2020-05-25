package main

import (
	"kanban-app-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	routes.Routes(router)

	log.Fatal(router.Run(":8080"))
}
