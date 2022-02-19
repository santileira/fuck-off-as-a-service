package service

import (
	"encoding/json"
	"fmt"
	"github.com/santileira/fuck-off-as-a-service/domain/message/domain"
	"github.com/santileira/fuck-off-as-a-service/http"
	"github.com/sirupsen/logrus"
)

const (
	foaasProtocol = "https"
	foaasDomain   = "foaas.com"
)

type MessageServiceImpl struct {
	FoaasProtocol string
	FoaasDomain   string
	client        http.Client
}

func NewMessageServiceImpl(client http.Client) *MessageServiceImpl {
	return &MessageServiceImpl{
		FoaasProtocol: foaasProtocol,
		FoaasDomain:   foaasDomain,
		client:        client,
	}
}

func (m *MessageServiceImpl) GetMessage(userID string) (*domain.Response, error) {
	url := fmt.Sprintf("%s://%s/thinking/%s/foaasAPI", m.FoaasProtocol, m.FoaasDomain, userID)
	body, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}

	response := &domain.Response{}
	if err := json.Unmarshal(body, response); err != nil {
		logrus.Errorf("Error unmarshaling the response, err: %s", err.Error())
		return nil, fmt.Errorf("error unmarshaling the body, err: %s", err.Error())
	}

	return response, nil
}
