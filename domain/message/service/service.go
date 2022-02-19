package service

import "github.com/santileira/fuck-off-as-a-service/domain/message/domain"

type MessageService interface {
	GetMessage(userID string) (*domain.Response, error)
}
