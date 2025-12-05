package main

import (
	"fmt"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/server"
)

func main() {

	limiter.StartLimiter()

	server.StartAPI()

	err := limiter.CloseDB()

	if err != nil {
		fmt.Errorf("Failed to close DB: %w", err)
	}

}
