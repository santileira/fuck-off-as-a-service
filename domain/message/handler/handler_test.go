package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/santileira/fuck-off-as-a-service/domain/message/domain"
	"github.com/santileira/fuck-off-as-a-service/domain/message/service"
	servicemocks "github.com/santileira/fuck-off-as-a-service/domain/message/service/mocks"
	"github.com/santileira/fuck-off-as-a-service/domain/message/validator"
	validatormocks "github.com/santileira/fuck-off-as-a-service/domain/message/validator/mocks"
	customhttp "github.com/santileira/fuck-off-as-a-service/http"
	"github.com/santileira/fuck-off-as-a-service/ratelimiter"
	ratelimitermocks "github.com/santileira/fuck-off-as-a-service/ratelimiter/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandleGetMessage(t *testing.T) {
	cases := []struct {
		name                 string
		userID               string
		mockMessageValidator *validatormocks.MessageValidator
		mockRateLimiter      *ratelimitermocks.RateLimiter
		mockMessageService   *servicemocks.MessageService
		expectedStatusCode   int
		expectedBody         string
	}{
		{
			"Should return an error when validator returns an error",
			"",
			func() *validatormocks.MessageValidator {
				mock := &validatormocks.MessageValidator{}
				mock.On("ValidateMessage", "").
					Return(fmt.Errorf("an error"))
				return mock
			}(),
			nil,
			nil,
			http.StatusBadRequest,
			`{"error":"an error"}`,
		},
		{
			"Should return an error when rate limiter returns false",
			"userID",
			func() *validatormocks.MessageValidator {
				mock := &validatormocks.MessageValidator{}
				mock.On("ValidateMessage", "userID").
					Return(nil)
				return mock
			}(),
			func() *ratelimitermocks.RateLimiter {
				mock := &ratelimitermocks.RateLimiter{}
				mock.On("AllowRequest", "userID").
					Return(false)
				return mock
			}(),
			nil,
			http.StatusTooManyRequests,
			`{"error":"request is not allowed now"}`,
		},
		{
			"Should return an error when service returns an error",
			"userID",
			func() *validatormocks.MessageValidator {
				mock := &validatormocks.MessageValidator{}
				mock.On("ValidateMessage", "userID").
					Return(nil)
				return mock
			}(),
			func() *ratelimitermocks.RateLimiter {
				mock := &ratelimitermocks.RateLimiter{}
				mock.On("AllowRequest", "userID").
					Return(true)
				return mock
			}(),
			func() *servicemocks.MessageService {
				mock := &servicemocks.MessageService{}
				mock.On("GetMessage", "userID").
					Return(nil, fmt.Errorf("error getting message"))
				return mock
			}(),
			http.StatusInternalServerError,
			`{"error":"error getting message"}`,
		},
		{
			"Should return a nil error",
			"userID",
			func() *validatormocks.MessageValidator {
				mock := &validatormocks.MessageValidator{}
				mock.On("ValidateMessage", "userID").
					Return(nil)
				return mock
			}(),
			func() *ratelimitermocks.RateLimiter {
				mock := &ratelimitermocks.RateLimiter{}
				mock.On("AllowRequest", "userID").
					Return(true)
				return mock
			}(),
			func() *servicemocks.MessageService {
				mock := &servicemocks.MessageService{}
				mock.On("GetMessage", "userID").
					Return(&domain.Response{
						Message:  "message",
						Subtitle: "subtitle",
					}, nil)
				return mock
			}(),
			http.StatusOK,
			`{"message":"message","subtitle":"subtitle"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			handler := NewMessageHandler(c.mockMessageValidator, c.mockRateLimiter, c.mockMessageService)

			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)
			context.Request, _ = http.NewRequest("GET", "/", nil)
			context.Request.Header.Set("User-Id", c.userID)

			// Operation
			handler.HandleGetMessage(context)

			// Validation
			assert.EqualValues(t, c.expectedStatusCode, w.Code)
			assert.EqualValues(t, c.expectedBody, w.Body.String())

			if c.mockMessageValidator != nil {
				c.mockMessageValidator.AssertNumberOfCalls(t, "ValidateMessage", 1)
			}

			if c.mockMessageService != nil {
				c.mockMessageService.AssertNumberOfCalls(t, "GetMessage", 1)
			}

			if c.mockRateLimiter != nil {
				c.mockRateLimiter.AssertNumberOfCalls(t, "AllowRequest", 1)
			}
		})
	}
}

func TestHandleGetMessageIntegration(t *testing.T) {
	// Initialization
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"message": "123, what the fuck were you actually thinking?","subtitle": "- foaasAPI"}`)
	}))
	defer server.Close()

	slidingWindowLogRateLimiter := ratelimiter.NewSlidingWindowLogRateLimiter(5, time.Millisecond*time.Duration(10000))
	httpClient := customhttp.NewClientImpl(time.Duration(5) * time.Second)
	messageService := service.NewMessageServiceImpl(httpClient)
	messageService.FoaasProtocol = "http"
	messageService.FoaasDomain = strings.Split(server.URL, "//")[1] // ex: http://127.0.0.1:61324
	messageValidator := validator.NewMessageValidatorImpl()
	messageHandler := NewMessageHandler(messageValidator, slidingWindowLogRateLimiter, messageService)

	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("User-Id", "UserIDSantiLeira")

	// Operation
	messageHandler.HandleGetMessage(context)

	// Validation
	assert.EqualValues(t, 200, w.Code)
	assert.EqualValues(t, `{"message":"123, what the fuck were you actually thinking?","subtitle":"- foaasAPI"}`,
		w.Body.String())
}
