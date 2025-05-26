package templates

func FrameworkInit(framework string) string {
	var mainContent string
	switch framework {
	case "Gin Gonic":
		mainContent = `
package main
import (
 	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World from Gin! 🌿",
		})
	})
	r.Run(":8080") // Listen and serve
}
`
		break
	case "Fiber":
		mainContent = `package main

import (
	"github.com/gofiber/fiber/v2"
)
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World from Fiber! ⚡",
		})
	})

	app.Listen(":8080")
}
`
		break
	case "Echo":
		mainContent = `package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello World from Echo! 🪶",
		})
	})

	e.Start(":8080")
}
`
	}
	return mainContent
}
