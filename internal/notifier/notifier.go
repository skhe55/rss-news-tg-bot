package notifier

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"news-feed-bot/internal/botkit/markup"
	"news-feed-bot/internal/model"
	"regexp"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkPosted(ctx context.Context, article model.Article) error
}

type Notifier struct {
	articles         ArticleProvider
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelId        int64
}

func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func New(articleProvider ArticleProvider, bot *tgbotapi.BotAPI, sendInterval time.Duration, lookupTimeWindow time.Duration, channelId int64) *Notifier {
	return &Notifier{
		articles:         articleProvider,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelId:        channelId,
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneArticles, err := n.articles.AllNotPosted(ctx, time.Now().Add(-n.lookupTimeWindow), 1)

	if err != nil {
		return err
	}

	if len(topOneArticles) == 0 {
		return nil
	}

	article := topOneArticles[0]
	summary, err := n.extractSummary(article)
	if err != nil {
		log.Printf("[ERROR] failed to extract summary: %v", err)
	}

	if err := n.sendArticle(article, summary); err != nil {
		return err
	}

	return n.articles.MarkPosted(ctx, article)

}

func (n *Notifier) sendArticle(article model.Article, summary string) error {
	const msgFormat = "*%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(n.channelId, fmt.Sprintf(
		msgFormat,
		markup.EscapeForMarkdown(article.Title),
		markup.EscapeForMarkdown(""),
		markup.EscapeForMarkdown(article.Link),
	))

	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (n *Notifier) extractSummary(article model.Article) (string, error) {
	var r io.Reader
	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}

	return "\n\n" + cleanText(doc.TextContent), nil
}

var redundantNewLines = regexp.MustCompile("\n{3}")

func cleanText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}
