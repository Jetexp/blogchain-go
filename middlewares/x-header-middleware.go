package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/zikwall/blogchain/types"
)

func XHeader(c *fiber.Ctx) {
	xHeader := types.NewXHeader(c)

	if xHeader.IsBlogchainApp() == false {
		c.Status(403).JSON(fiber.Map{
			"status":  100,
			"message": "Only real blogchain clients can use the API",
		})

		return
	}

	requestPath := c.Path()
	fmt.Println(fmt.Sprintf("Request %s from platform: %s@%s",
		requestPath,
		xHeader.XPlatform,
		xHeader.XAppVersion,
	))

	c.Next()
}