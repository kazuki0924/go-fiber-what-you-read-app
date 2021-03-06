package infrastructure

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func timer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// start timer
		start := time.Now()
		// next routes
		err := c.Next()
		// stop timer
		stop := time.Now()

		c.Append("Server-Timing", fmt.Sprintf("app;dur=%v", stop.Sub(start).String()))
		return err
	}
}

var cacheConfig = cache.Config{
	Next: func(c *fiber.Ctx) bool {
		return c.Query("refresh") == "true"
	},
	Expiration:   15 * time.Minute,
	CacheControl: true,
}

var loggerConfig = logger.Config{
	Format: "${cyan}[${time}] ${white}${pid} ${red}${status} ${blue}[${method}] ${white}${path}\n",
}

func SetupFiberMiddleWares(app *fiber.App) *fiber.App {
	app.Use(
		timer(),
		logger.New(loggerConfig),
		cors.New(cors.Config{
			Next:             nil,
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
			MaxAge:           0,
		}),
		cache.New(cacheConfig),
	)
	return app
}
