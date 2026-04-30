package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/xoejang/gitel/internal/config"
	"github.com/xoejang/gitel/internal/handler"
	"github.com/xoejang/gitel/internal/service"
	"github.com/xoejang/gitel/pkg/llm"
	"github.com/xoejang/gitel/pkg/telegram"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	extractor := service.NewExtractor()
	llmClient := llm.NewClient(cfg.LLM.APIKey, cfg.LLM.Model, cfg.LLM.BaseURL, cfg.LLM.Timeout)
	formatter := service.NewFormatter(llmClient)
	telegramClient := telegram.NewClient(cfg.Telegram.BotToken, cfg.Telegram.ChatID)
	processor := service.NewProcessor(formatter, telegramClient)

	webhookHandler := handler.NewWebhookHandler(cfg.Webhook.Secret, extractor, processor)
	server := handler.NewServer(cfg.Server.Port, webhookHandler)

	go func() {
		slog.Info("server starting", "port", cfg.Server.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down gracefully")
	if err := server.Shutdown(); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("waiting for background jobs to complete")
	processor.Wait()

	slog.Info("server exited")
}
