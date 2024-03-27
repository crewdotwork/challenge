package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type enrichmentRequest struct {
	LinkedInProfile string `json:"linkedin_profile"`
}

type enrichmentResponse struct {
	Emails []string `json:"emails"`
	Phones []string `json:"phones"`
}

var enrichments = map[string]enrichmentResponse{
	"https://www.linkedin.com/in/steve-jobs": {
		Emails: []string{"steve@apple.com", "stevejobs@gmail.com"},
		Phones: []string{"+14155552671", "+14155552672"},
	},
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)

	go func() {
		<-done
		log.Debug().Msg("received signal to shutdown server")
		cancel()
	}()

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}))

	router.POST("/api/v1/enrich", func(c *gin.Context) {
		var req enrichmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(req.LinkedInProfile)
		if req.LinkedInProfile == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "linkedin_profile is required"})
			return
		}

		enrichment, exists := enrichments[req.LinkedInProfile]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "enrichment not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"enrichment": enrichment})
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("failed to start server")
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		<-ctx.Done()
		if err := srv.Shutdown(context.WithoutCancel(ctx)); err != nil {
			log.Error().Err(err).Msg("failed to shutdown server")
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		log.Error().Err(err).Msg("failed to run errgorup")
	}
}
