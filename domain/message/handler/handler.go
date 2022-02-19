package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/santileira/fuck-off-as-a-service/domain/message/service"
	"github.com/santileira/fuck-off-as-a-service/domain/message/validator"
	"github.com/santileira/fuck-off-as-a-service/ratelimiter"
	"github.com/sirupsen/logrus"
	"net/http"
)

const userIDHeader = "User-Id"

type MessageHandler struct {
	messageValidator validator.MessageValidator
	rateLimiter      ratelimiter.RateLimiter
	messageService   service.MessageService
}

func NewMessageHandler(messageValidator validator.MessageValidator,
	rateLimiter ratelimiter.RateLimiter,
	messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageValidator: messageValidator,
		rateLimiter:      rateLimiter,
		messageService:   messageService,
	}
}

func (m *MessageHandler) HandleGetMessage(ginContext *gin.Context) {
	userID := ginContext.GetHeader(userIDHeader)
	if err := m.messageValidator.ValidateMessage(userID); err != nil {
		logrus.Errorf("Error validating the message, userID: %s", userID)
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !m.rateLimiter.AllowRequest(userID) {
		logrus.Errorf("Request is not allowed for user: %s", userID)
		ginContext.JSON(http.StatusTooManyRequests, gin.H{
			"error": "request is not allowed now",
		})
		return
	}

	response, err := m.messageService.GetMessage(userID)
	if err != nil {
		logrus.Errorf("Error getting the message, userID: %s, err: %s", userID, err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ginContext.JSON(http.StatusOK, response)
	return
}
