package bot

import (
	"context"
	"fmt"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdAddSource(storage SourceStorage) botkit.ViewFunc {
	type addSourceArgs struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			if _, sendError := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Возникла ошибка при парсинге json'a: %v. Проверьте что ввод корректен. \n Пример: {'name': 'type_name', 'url': 'type_url'}", err))); sendError != nil {
				return sendError
			}
			return err
		}

		source := model.Source{
			Name:    args.Name,
			FeedUrl: args.URL,
		}

		sourceId, err := storage.Add(ctx, source)
		if err != nil {
			return nil
		}

		var (
			msgText = fmt.Sprintf(
				"Источник добавлен с ID: '%d'\\. Используйте этот ID для управления источником\\.",
				sourceId,
			)
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
