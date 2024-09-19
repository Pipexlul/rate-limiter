package interfaces

import (
	"net/http"

	"github.com/pipexlul/rate-limiter/models"
)

type Limiter interface {
	SetPolicy(endpoint string, policy models.Policy)
	SetPolicyForEndpoints(endpoints []string, policy models.Policy)
	Middleware(next http.Handler) http.Handler
}
