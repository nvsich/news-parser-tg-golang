package notifier

import (
	"context"
	"fmt"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"news-parser-tg/internal/bot/markup"
	"news-parser-tg/internal/model"
	"strings"
	"time"
)

type ArticleProvider interface {
	AllNonPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkPosted(ctx context.Context, id int64) error
}

type Summarizer interface {
	Summarize(ctx context.Context, text string) (string, error)
}

type Notifier struct {
	articleProvider  ArticleProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelID        int64
}

func New(
	articleProvider ArticleProvider,
	summarizer Summarizer,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookupTimeWindow time.Duration,
	channelID int64,
) *Notifier {
	return &Notifier{
		articleProvider:  articleProvider,
		summarizer:       summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelID:        channelID,
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneArticle, err := n.articleProvider.AllNonPosted(ctx, time.Now().Add(-n.lookupTimeWindow), 1)
	if err != nil {
		return err
	}

	if len(topOneArticle) == 0 {
		return nil
	}

	article := topOneArticle[0]

	summary, err := n.extractSummary(ctx, &article)
	if err != nil {
		return err
	}

	if err := n.sendArticle(&article, summary); err != nil {
		return err
	}

	return n.articleProvider.MarkPosted(ctx, article.ID)
}

func (n *Notifier) extractSummary(ctx context.Context, article *model.Article) (string, error) {
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

	summary, err := n.summarizer.Summarize(ctx, doc.TextContent)

	return "\n\n" + summary, nil
}

func (n *Notifier) sendArticle(article *model.Article, summary string) error {
	const msgFormat = "*%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(
		n.channelID,
		fmt.Sprintf(
			msgFormat,
			markup.EscapeForMarkdown(article.Title),
			markup.EscapeForMarkdown(summary),
			markup.EscapeForMarkdown(article.Link),
		),
	)

	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
