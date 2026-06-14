package utils

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(component templ.Component) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Status(http.StatusOK)
		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
