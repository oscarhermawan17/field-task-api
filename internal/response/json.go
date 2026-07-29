package response

import "github.com/gin-gonic/gin"

// JSON centralizes JSON responses so handlers keep their focus on HTTP behavior.
func JSON(c *gin.Context, status int, body any) {
	c.JSON(status, body)
}
