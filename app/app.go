package app

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

var instance *App
var once sync.Once

func GetInstance() *App {
	once.Do(func() {
		r := gin.Default()
		r.Use(securityMiddleware())
		r.MaxMultipartMemory = 1 << 20 // 1 MB multipart limit
		instance = &App{Router: r}
	})
	return instance
}

func securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			// Immutable assets — GLBs, CSS, JS never change without a new deploy
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Pages: allow caching for 5 minutes, revalidate after
			c.Header("Cache-Control", "public, max-age=300, must-revalidate")
		}

		c.Next()
	}
}
