package bot

import (
	"context"
	"fmt"
	"news-feed-bot/internal/botkit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdDeleteSource(storage SourceStorage) botkit.ViewFunc {
	type deleteSourceArgs struct {
		Id int64 `json:"id"`
	}
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[deleteSourceArgs](update.Message.CommandArguments())

		if err != nil {
			if _, sendError := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Возникла ошибка при парсинге json'a: %v\\. Проверьте что ввод корректен\\.", err))); sendError != nil {
				return sendError
			}
			return err
		}

		sourceId := args.Id

		if err := storage.Delete(ctx, sourceId); err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf("Источник с ID:%d успешно удалён из базы данных\\.", sourceId)
			reply   = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
