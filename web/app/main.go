package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/images", "web/assets/images")

	setupRoutes(e)

	if err := e.Start(":6969"); err != nil {
		log.Fatal(err)
	}

	log.Default().Println("Running templ server on 9090")
}
