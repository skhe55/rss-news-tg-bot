package bot

import (
	"context"
	"fmt"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/botkit/markup"
	"news-feed-bot/internal/model"
	"news-feed-bot/internal/utils"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdListSources(storage SourceStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := storage.Sources(ctx)

		if err != nil {
			return err
		}

		var (
			sourceInfos = utils.Map(sources, func(source model.Source, _ int) string {
				return formatSource(source)
			})
			msgText = fmt.Sprintf(
				"Список источников \\(всего %d\\):\n\n%s",
				len(sources),
				strings.Join(sourceInfos, "\n\n"),
			)
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}

func formatSource(source model.Source) string {
	return fmt.Sprintf(
		"*%s*\nID: `%d`\nURL фида: %s",
		markup.EscapeForMarkdown(source.Name),
		source.Id,
		markup.EscapeForMarkdown(source.FeedUrl),
	)
}
