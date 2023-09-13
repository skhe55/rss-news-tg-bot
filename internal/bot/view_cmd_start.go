package bot

import (
	"context"
	"news-feed-bot/internal/botkit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Привет, вот список моих команд:\n/addsource\n/listsources\n/deletesource\n/updatesource\n")); err != nil {
			return err
		}
		return nil
	}
}
