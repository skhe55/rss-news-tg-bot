package bot

import (
	"context"
	"fmt"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdUpdateSource(storage SourceStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[model.UpdateSource](update.Message.CommandArguments())
		if err != nil {
			if _, sendError := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Возникла ошибка при парсинге json'a: \n%v.\n Проверьте что ввод корректен. \n Пример: {'id': 999, 'name': 'type_name', 'url': 'type_url'}\n id - обязателен.\n", err))); sendError != nil {
				return sendError
			}
			return err
		}

		source := model.UpdateSource{
			Id:   args.Id,
			Name: args.Name,
			URL:  args.URL,
		}

		if err := storage.Update(ctx, source); err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf(
				"Источник ID: '%d'\\ успешно обновлён",
				args.Id,
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
