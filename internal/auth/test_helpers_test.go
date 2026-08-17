package auth

import "github.com/gin-gonic/gin"

func newTestRouter(handler *Handler, secret []byte) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/login", handler.Login)
	router.GET("/auth/me", RequireAuth(secret), handler.Me)
	return router
}
