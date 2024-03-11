package middlewares

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

// Throttler represents the throttling middleware.
type Throttler struct {
	mu         sync.Mutex         // Mutex to synchronize access to the clients map
	clients    map[string]*client // Map to store client rate limiters and last seen times
	limit      rate.Limit         // Maximum number of requests allowed per second
	burstLimit int                // Maximum number of requests allowed to burst per second
}

// client holds the rate limiter and last seen time for each client.
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewThrottler creates a new Throttler instance.
func NewThrottler(limit rate.Limit, burstLimit int) *Throttler {
	return &Throttler{
		clients:    make(map[string]*client),
		limit:      limit,
		burstLimit: burstLimit,
	}
}

// Middleware restricts request rate per client IP address.
func (t *Throttler) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := getClientIP(c.Request())

		t.mu.Lock()
		defer t.mu.Unlock()

		// Create a new client entry if not exists
		if _, found := t.clients[ip]; !found {
			t.clients[ip] = &client{
				limiter: rate.NewLimiter(t.limit, t.burstLimit),
			}
		}

		// Update last seen time for the client
		t.clients[ip].lastSeen = time.Now()

		// Check if client has exceeded rate limit
		if !t.clients[ip].limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, echo.Map{
				"status": "Request Failed",
				"body":   "Too many requests, try again later.",
			})
		}

		return next(c)
	}
}

// getClientIP extracts the client IP address from the request.
func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}
