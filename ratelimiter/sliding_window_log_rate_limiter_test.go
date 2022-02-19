package ratelimiter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetRequestsInTheWindowTime(t *testing.T) {
	cases := []struct {
		name                          string
		inputRequests                 []time.Time
		inputNow                      time.Time
		inputRateWindowInMilliseconds int
		expectedOutput                []time.Time
	}{
		{
			"Should return empty when requests are nil",
			nil,
			time.Now(),
			10,
			[]time.Time{},
		},
		{
			"Should return empty when requests are empty",
			[]time.Time{},
			time.Now(),
			10,
			[]time.Time{},
		},
		{
			"Should return empty when all the requests were on the past",
			[]time.Time{
				time.Date(2022, time.January, 18, 16, 10, 10, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 10, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 11, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 12, 00, time.UTC),
			},
			time.Date(2022, time.January, 18, 16, 11, 15, 00, time.UTC),
			1000,
			[]time.Time{},
		},
		{
			"Should return the last 2 requests when they are inside of the window time, " +
				"now is: 2022-02-18 16:10:15 00" +
				"past two requests: 2022-02-18 16:10:11, 2022-02-18 16:10:12, 2022-02-18 16:10:14",
			[]time.Time{
				time.Date(2022, time.January, 18, 16, 10, 10, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 11, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 12, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 14, 00, time.UTC),
			},
			time.Date(2022, time.January, 18, 16, 10, 15, 00, time.UTC),
			1000 * 4,
			[]time.Time{
				time.Date(2022, time.January, 18, 16, 10, 11, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 12, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 14, 00, time.UTC),
			},
		},
		{
			"Should return the all the requests because they are inside the window time",
			[]time.Time{
				time.Date(2022, time.January, 18, 16, 10, 05, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 11, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 15, 00, time.UTC),
			},
			time.Date(2022, time.January, 18, 16, 10, 15, 00, time.UTC),
			1000 * 10,
			[]time.Time{
				time.Date(2022, time.January, 18, 16, 10, 05, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 11, 00, time.UTC),
				time.Date(2022, time.January, 18, 16, 10, 15, 00, time.UTC),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			rateLimiter := NewSlidingWindowLogRateLimiter(0,
				time.Duration(c.inputRateWindowInMilliseconds)*time.Millisecond)

			// Operation
			requests := rateLimiter.getRequestsInTheWindowTime(c.inputRequests, c.inputNow)

			// Validation
			assert.EqualValues(t, c.expectedOutput, requests)
		})
	}
}

func TestAllowRequestShouldReturnTrueWhenItIsTheFirstRequest(t *testing.T) {
	// Initialization
	userID := "userID"

	slidingWindowLogRateLimiter := NewSlidingWindowLogRateLimiter(5, time.Duration(10000)*time.Millisecond)
	slidingWindowLogRateLimiter.now = func() time.Time {
		return time.Date(2022, time.January, 10, 10, 10, 10, 00, time.UTC)
	}

	// Operation
	isAllowed := slidingWindowLogRateLimiter.AllowRequest(userID)

	// Validation
	assert.True(t, isAllowed)
	assert.Len(t, slidingWindowLogRateLimiter.requestsByUser, 1)
	assert.EqualValues(t, slidingWindowLogRateLimiter.requestsByUser["userID"],
		[]time.Time{slidingWindowLogRateLimiter.now()})
}

func TestAllowRequestShouldReturnTrueWhenRequestsInWindowTimeIsLessThanLimit(t *testing.T) {
	// Initialization
	userID := "userID"

	slidingWindowLogRateLimiter := NewSlidingWindowLogRateLimiter(5, time.Duration(10000)*time.Millisecond)
	slidingWindowLogRateLimiter.now = func() time.Time {
		return time.Date(2022, time.January, 10, 10, 10, 18, 00, time.UTC)
	}

	slidingWindowLogRateLimiter.requestsByUser["userID"] = []time.Time{
		time.Date(2022, time.January, 10, 10, 10, 01, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 02, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 13, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 14, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 17, 00, time.UTC),
	}

	// Operation
	isAllowed := slidingWindowLogRateLimiter.AllowRequest(userID)

	// Validation
	assert.True(t, isAllowed)
	assert.Len(t, slidingWindowLogRateLimiter.requestsByUser, 1)
	assert.EqualValues(t, slidingWindowLogRateLimiter.requestsByUser["userID"],
		[]time.Time{
			time.Date(2022, time.January, 10, 10, 10, 13, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 14, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 17, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 18, 00, time.UTC),
		})
}

func TestAllowRequestShouldReturnFalseWhenRequestsInWindowTimeIsEqualThanLimit(t *testing.T) {
	// Initialization
	userID := "userID"

	slidingWindowLogRateLimiter := NewSlidingWindowLogRateLimiter(3, time.Duration(10000)*time.Millisecond)
	slidingWindowLogRateLimiter.now = func() time.Time {
		return time.Date(2022, time.January, 10, 10, 10, 18, 00, time.UTC)
	}

	slidingWindowLogRateLimiter.requestsByUser["userID"] = []time.Time{
		time.Date(2022, time.January, 10, 10, 10, 01, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 02, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 13, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 14, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 17, 00, time.UTC),
	}

	// Operation
	isAllowed := slidingWindowLogRateLimiter.AllowRequest(userID)

	// Validation
	fmt.Println(slidingWindowLogRateLimiter.requestsByUser["userID"])
	assert.False(t, isAllowed)
	assert.Len(t, slidingWindowLogRateLimiter.requestsByUser, 1)
	assert.EqualValues(t, slidingWindowLogRateLimiter.requestsByUser["userID"],
		[]time.Time{
			time.Date(2022, time.January, 10, 10, 10, 13, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 14, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 17, 00, time.UTC),
		})
}

func TestAllowRequestShouldReturnFalseWhenRequestsInWindowTimeIsGreaterThanLimit(t *testing.T) {
	// Initialization
	userID := "userID"

	slidingWindowLogRateLimiter := NewSlidingWindowLogRateLimiter(3, time.Duration(10000)*time.Millisecond)
	slidingWindowLogRateLimiter.now = func() time.Time {
		return time.Date(2022, time.January, 10, 10, 10, 18, 00, time.UTC)
	}

	slidingWindowLogRateLimiter.requestsByUser["userID"] = []time.Time{
		time.Date(2022, time.January, 10, 10, 10, 01, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 02, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 13, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 14, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 15, 00, time.UTC),
		time.Date(2022, time.January, 10, 10, 10, 17, 00, time.UTC),
	}

	// Operation
	isAllowed := slidingWindowLogRateLimiter.AllowRequest(userID)

	// Validation
	fmt.Println(slidingWindowLogRateLimiter.requestsByUser["userID"])
	assert.False(t, isAllowed)
	assert.Len(t, slidingWindowLogRateLimiter.requestsByUser, 1)
	assert.EqualValues(t, slidingWindowLogRateLimiter.requestsByUser["userID"],
		[]time.Time{
			time.Date(2022, time.January, 10, 10, 10, 13, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 14, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 15, 00, time.UTC),
			time.Date(2022, time.January, 10, 10, 10, 17, 00, time.UTC),
		})
}
