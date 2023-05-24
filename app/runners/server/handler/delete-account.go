package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *RouteHandler) DeleteAccountHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "account successfully deleted",
	})
}
