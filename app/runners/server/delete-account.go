package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) DeleteAccountHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "account successfully deleted",
	})
}
