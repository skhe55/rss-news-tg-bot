package middleware

import (
	"context"
	"news-feed-bot/internal/botkit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AdminOnly(channelId int64, next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		admines, err := bot.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelId,
				},
			},
		)

		if err != nil {
			return err
		}

		for _, admin := range admines {
			if admin.User.ID == update.Message.From.ID {
				return next(ctx, bot, update)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав для выполнения данной команды.")); err != nil {
			return err
		}

		return nil
	}
}
