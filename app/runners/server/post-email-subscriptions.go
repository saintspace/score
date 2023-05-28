package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostEmailSubscriptionRequestData struct {
	EmailAddress string `json:"email"`
}

func (s *Server) PostEmailSubscriptionsHandler(c *gin.Context) {
	var data PostEmailSubscriptionRequestData
	if err := c.ShouldBindJSON(&data); err != nil {
		s.logger.ErrorWithContext(
			"error while parsing PostEmailSubscriptionRequestData",
			"error", err.Error(),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !s.emailService.IsValidEmail(data.EmailAddress) {
		s.logger.ErrorWithContext(
			"invalid email address",
			"email", data.EmailAddress,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email address"})
		return
	}
	subscriptionExists, err := s.emailService.EmailSubscriptionExists(data.EmailAddress)
	if err != nil {
		s.logger.ErrorWithContext(
			"error while checking if email subscription already exists",
			"error", err.Error(),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while processing your email subscription"})
		return
	}

	if !subscriptionExists {
		if err := s.emailService.CreateEmailSubscription(data.EmailAddress); err != nil {
			s.logger.ErrorWithContext(
				"error while creating email subscription",
				"error", err.Error(),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while saving your email subscription"})
			return
		} else {
			s.logger.InfoWithContext(
				"Successfully created email subscription",
				"email", data.EmailAddress,
			)
		}
	} else {
		s.logger.InfoWithContext(
			"Email subscription already exists",
			"email", data.EmailAddress,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "email subscription successful",
		"email":   data.EmailAddress,
	})
}
