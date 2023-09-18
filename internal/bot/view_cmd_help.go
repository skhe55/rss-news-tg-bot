package bot

import (
	"context"
	"news-feed-bot/internal/botkit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdHelp() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(
			tgbotapi.NewMessage(update.FromChat().ID,
				"Привет, я  умею работать только с загрузкой/редактированием/удалением источника, а так же я могу отдать список существующих источников!\\")); err != nil {
			return err
		}
		return nil
	}
}
