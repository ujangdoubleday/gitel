package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/xoejang/gitel/internal/config"
	"github.com/xoejang/gitel/internal/handler"
	"github.com/xoejang/gitel/internal/service"
	"github.com/xoejang/gitel/pkg/llm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	extractor := service.NewExtractor()
	llmClient := llm.NewClient(cfg.LLM.APIKey, cfg.LLM.Model, cfg.LLM.BaseURL, cfg.LLM.Timeout)
	formatter := service.NewFormatter(llmClient)
	webhookHandler := handler.NewWebhookHandler(cfg.Webhook.Secret, extractor, formatter)
	server := handler.NewServer(cfg.Server.Port, webhookHandler)

	go func() {
		log.Printf("server starting on port %s", cfg.Server.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down gracefully...")
	if err := server.Shutdown(); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
