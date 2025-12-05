package server

import (
	"log"
	"net/http"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/quotes"
	"time"

	"github.com/gin-gonic/gin"
)

func StartAPI() {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/get-quote", func(ctx *gin.Context) {

		ip := ctx.ClientIP()

		allowed, err := limiter.CanAccess(ip)

		if err != nil {
			log.Printf("Error checking rate limit: %w", err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		if allowed == false {
			var val time.Time

			val, err = limiter.GetOldestLog(ip)

			if err != nil {
				log.Printf("Failed to get 'try_in': %w", err)

				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}

			tryAgainIn := limiter.GetPeriod() - int(time.Since(val).Seconds())

			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"try_in": tryAgainIn,
			})

			return
		}

		quote, err := quotes.PickQuoteFromAuthor(ctx.Query("author"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		limiter.AddLog(ip)

		ctx.JSON(http.StatusOK, gin.H{
			"text":   quote.Text,
			"author": quote.Author,
		})
	})

	router.Run()

}
