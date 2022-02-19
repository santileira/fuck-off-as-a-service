package server

import (
	"github.com/santileira/fuck-off-as-a-service/domain/message/handler"
	"github.com/santileira/fuck-off-as-a-service/domain/message/service"
	"github.com/santileira/fuck-off-as-a-service/domain/message/validator"
	"github.com/santileira/fuck-off-as-a-service/http"
	"github.com/santileira/fuck-off-as-a-service/ratelimiter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

type Runnable struct{}

func NewRunnable() *Runnable {
	return &Runnable{}
}

func (r *Runnable) Cmd() *cobra.Command {
	options := &Options{}

	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Runs fuck off as a service server",
		Long:  `Runs fuck of as a service server`,
	}

	cmd.Flags().StringVar(&options.LogLevel, "log-level", defaultLogLevel, "log leve to use")
	cmd.Flags().IntVar(&options.RateLimitCount, "rate-limit-count", defaultRateLimitCount, "maximum quantity of requests "+
		"that a user can do in a window of time")
	cmd.Flags().IntVar(&options.RateLimitWindowInMilliseconds, "rate-limit-window-in-milliseconds", defaultRateLimitWindowInMilliseconds,
		"window of time in milliseconds to limit the quantity of requests that a user can do")
	cmd.Flags().IntVar(&options.TimeoutInMilliseconds, "timeout-in-milliseconds", defaultTimeoutInMilliseconds,
		"timeout of the api calls")

	cmd.Run = func(_ *cobra.Command, _ []string) {
		server := r.Run(options)
		server.Start()
	}
	return cmd
}

func (r *Runnable) Run(options *Options) *Server {
	r.configureLog(options.LogLevel)

	slidingWindowLogRateLimiter := ratelimiter.NewSlidingWindowLogRateLimiter(
		options.RateLimitCount,
		time.Duration(options.RateLimitWindowInMilliseconds)*time.Millisecond)

	httpClient := http.NewClientImpl(time.Duration(options.TimeoutInMilliseconds) * time.Millisecond)

	messageService := service.NewMessageServiceImpl(httpClient)
	messageValidator := validator.NewMessageValidatorImpl()
	messageHandler := handler.NewMessageHandler(messageValidator, slidingWindowLogRateLimiter, messageService)

	return NewServer(messageHandler)
}

func (r *Runnable) configureLog(logLevel string) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Warnf("Error passing the log level: %s", logLevel)
		return
	}

	logrus.Infof("Setting log level: %s", lvl.String())
	logrus.SetLevel(lvl)
}
