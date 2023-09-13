package main

import (
	"context"
	"errors"
	"log"
	"news-feed-bot/internal/bot"
	"news-feed-bot/internal/bot/middleware"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/config"
	"news-feed-bot/internal/fetcher"
	"news-feed-bot/internal/notifier"
	"news-feed-bot/internal/storage"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)

	if err != nil {
		log.Printf("failed to create bot: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcher        = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)

		notifier = notifier.New(
			articleStorage,
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelId,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart())
	newsBot.RegisterCmdView("addsource", middleware.AdminOnly(config.Get().TelegramChannelId, bot.ViewCmdAddSource(sourceStorage)))
	newsBot.RegisterCmdView("listsources", middleware.AdminOnly(config.Get().TelegramChannelId, bot.ViewCmdListSources(sourceStorage)))
	newsBot.RegisterCmdView("deletesource", middleware.AdminOnly(config.Get().TelegramChannelId, bot.ViewCmdDeleteSource(sourceStorage)))
	newsBot.RegisterCmdView("updatesources", middleware.AdminOnly(config.Get().TelegramChannelId, bot.ViewCmdUpdateSource(sourceStorage)))

	tgCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "addsource",
			Description: "Создать источник",
		},
		tgbotapi.BotCommand{
			Command:     "listsources",
			Description: "Получить список источников",
		},
		tgbotapi.BotCommand{
			Command:     "deletesource",
			Description: "Удалить источник",
		},
		tgbotapi.BotCommand{
			Command:     "updatesources",
			Description: "Обновить источник",
		},
	)

	_, _ = botAPI.Request(tgCfg)

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start fetcher: %v", err)
				return
			}

			log.Printf("fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start notifier: %v", err)
				return
			}

			log.Printf("notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to run bot: %v", err)
			return
		}
	}

	log.Printf("bot stopped")
}
