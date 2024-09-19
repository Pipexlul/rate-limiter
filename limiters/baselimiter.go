package limiters

import (
	"fmt"
	"github.com/pipexlul/rate-limiter/interfaces"
	"net/http"
	"sync"
	"time"

	"github.com/pipexlul/rate-limiter/models"
)

type (
	policyKey = string
	clientKey = string
)

var _ interfaces.Limiter = (*BaseRateLimiter)(nil)

type BaseRateLimiter struct {
	defaultPolicy models.Policy

	policies map[policyKey]models.Policy
	clients  map[clientKey]*models.ClientData

	mu sync.Mutex
}

func (rl *BaseRateLimiter) SetPolicy(endpoint string, policy models.Policy) {
	rl.policies[endpoint] = policy
}

func (rl *BaseRateLimiter) SetPolicyForEndpoints(endpoints []string, policy models.Policy) {
	for _, endpoint := range endpoints {
		rl.policies[endpoint] = policy
	}
}

func (rl *BaseRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allowClient(r.RemoteAddr, r.URL.Path) {
			http.Error(w, "Too many requests!", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewBaseRateLimiter(defaultPolicy models.Policy) *BaseRateLimiter {
	return &BaseRateLimiter{
		defaultPolicy: defaultPolicy,
		policies:      make(map[policyKey]models.Policy),
		clients:       make(map[clientKey]*models.ClientData),
	}
}

func (rl *BaseRateLimiter) allowClient(clientIP, endpoint string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	policy := rl.defaultPolicy
	if policyOverride, exists := rl.policies[endpoint]; exists {
		policy = policyOverride
	}

	clientKey := fmt.Sprintf("%s_%s", endpoint, clientIP)
	client, exists := rl.clients[clientKey]
	if !exists || time.Since(client.LastRequest) > policy.Interval {
		rl.clients[clientKey] = &models.ClientData{
			Tokens:      policy.MaxRequests - 1,
			LastRequest: time.Now(),
		}
		return true
	}

	if client.Tokens > 0 {
		client.Tokens--
		client.LastRequest = time.Now()
		return true
	}

	return false
}
