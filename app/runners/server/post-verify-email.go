package server

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type PostVerifyEmailRequestData struct {
	SubscriptionToken string `json:"token"`
}

func (s *Server) PostVerifyEmailHandler(c *gin.Context) {
	var data PostVerifyEmailRequestData
	if err := c.ShouldBindJSON(&data); err != nil {
		s.logger.ErrorWithContext(
			"error while parsing PostVerifyEmailHandler request",
			"error", err.Error(),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	unescapedSubscriptionToken, err := url.QueryUnescape(data.SubscriptionToken)
	if err != nil {
		s.logger.ErrorWithContext(
			"error while unescaping email verification token",
			"error", err.Error(),
			"token", data.SubscriptionToken,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while parsing verification token"})
		return
	}
	err = s.emailService.VerifyEmailWithSubscriptionToken(unescapedSubscriptionToken)
	if err != nil {
		s.logger.ErrorWithContext(
			"error while verifying email",
			"error", err.Error(),
			"token", unescapedSubscriptionToken,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error while attempting to verify subscription",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "email verification successful",
	})
}
