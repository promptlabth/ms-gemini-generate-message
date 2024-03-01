// build go1.16

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	promptMessage "github.com/promptlabth/ms-gemini-generate-message/app/prompt_message"
	"github.com/promptlabth/ms-gemini-generate-message/database"
	gemini "github.com/promptlabth/ms-gemini-generate-message/services"
)

func main() {

	db, err := database.Open(
		database.Config{
			UserName:     os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASS"),
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			DatabaseName: os.Getenv("DB_NAME"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := NewRouter(db)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func NewRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.Use(
		func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		},
	)

	gen := gemini.NewGeminiService()

	promptMessageStorage := promptMessage.NewPromptMessageStorage(db)
	promptMessageHandler := promptMessage.NewPromptMessageHandler(*promptMessageStorage, gen)
	router.POST("/generate-message", promptMessageHandler.GenerateMessage)

	router.GET("/stream", promptMessageHandler.GenerateStream)
	return router
}
