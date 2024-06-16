package server

import (
	"fmt"
	"net/http"

	"github.com/ryancarlos88/go-rate-limiter/config"
	ratelimiter "github.com/ryancarlos88/go-rate-limiter/internal/domain/rate-limiter"
	helloworld "github.com/ryancarlos88/go-rate-limiter/internal/web/handler/hello-world"
	"github.com/ryancarlos88/go-rate-limiter/internal/web/middleware/limiter"
)

type Server struct {
	RateLimiter *ratelimiter.RateLimiter
	Cfg         *config.Config
}

func NewServer(cfg *config.Config) *Server {
	rl := ratelimiter.NewRateLimiter(cfg)
	return &Server{
		RateLimiter: rl,
		Cfg:         cfg,
	}
}

func (s *Server) Start() {
	http.Handle("/", limiter.LimitMiddleware(http.HandlerFunc(helloworld.HelloWorldHandler), s.RateLimiter))

	http.ListenAndServe(":8080", nil)
}
func init() {
	fmt.Println("Server starting...")
}
