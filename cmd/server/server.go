package server

import (
	"github.com/gin-gonic/gin"
	"github.com/santileira/fuck-off-as-a-service/domain/message/handler"
)

type Server struct {
	messageHandler *handler.MessageHandler
}

func NewServer(messageHandler *handler.MessageHandler) *Server {
	return &Server{
		messageHandler: messageHandler,
	}
}

func (s *Server) Start() {
	engine := gin.Default()
	s.attachEndpoints(engine)
	if err := engine.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) attachEndpoints(engine *gin.Engine) {
	engine.GET("/message", s.messageHandler.HandleGetMessage)
}
