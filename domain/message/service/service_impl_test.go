package service

import (
	"fmt"
	"github.com/santileira/fuck-off-as-a-service/domain/message/domain"
	httpmock "github.com/santileira/fuck-off-as-a-service/http/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMessage(t *testing.T) {
	cases := []struct {
		name             string
		userID           string
		mockClient       *httpmock.Client
		expectedResponse *domain.Response
		expectedError    error
	}{
		{
			"Should return an error when client returns an error",
			"userID",
			func() *httpmock.Client {
				mock := &httpmock.Client{}
				mock.On("Get", "https://foaas.com/thinking/userID/foaasAPI").
					Return(nil, fmt.Errorf("error getting response from foaas"))
				return mock
			}(),
			nil,
			fmt.Errorf("error getting response from foaas"),
		},
		{
			"Should return an error when there's an error unmarshalling the response",
			"userID",
			func() *httpmock.Client {
				mock := &httpmock.Client{}
				mock.On("Get", "https://foaas.com/thinking/userID/foaasAPI").
					Return(nil, nil)
				return mock
			}(),
			nil,
			fmt.Errorf("error unmarshaling the body, err: unexpected end of JSON input"),
		},
		{
			"Should return an error when there's an error unmarshalling the response",
			"userID",
			func() *httpmock.Client {
				mock := &httpmock.Client{}
				mock.On("Get", "https://foaas.com/thinking/userID/foaasAPI").
					Return([]byte(`{"message": "userID, what the fuck were you actually thinking?","subtitle": "- foaasAPI"}`),
						nil)
				return mock
			}(),
			&domain.Response{Message: "userID, what the fuck were you actually thinking?",
				Subtitle: "- foaasAPI"},
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			service := NewMessageServiceImpl(c.mockClient)

			// Operation
			response, err := service.GetMessage(c.userID)

			// Validation
			assert.EqualValues(t, c.expectedResponse, response)
			assert.EqualValues(t, c.expectedError, err)
			c.mockClient.AssertNumberOfCalls(t, "Get", 1)
		})
	}
}
