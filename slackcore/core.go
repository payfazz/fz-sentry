package slackcore

import (
	"fmt"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type slackCore struct {
	zapcore.Core

	hookURL  string
	minLevel zapcore.Level

	fields []zapcore.Field
}

// NewWrapper ...
func NewWrapper(hookURL string, minLevel zapcore.Level) func(zapcore.Core) zapcore.Core {
	return func(c zapcore.Core) zapcore.Core {
		return &slackCore{
			Core:     c,
			hookURL:  hookURL,
			minLevel: minLevel,
		}
	}
}

func (c *slackCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	// Join core fields with entry fields
	fields = append(
		append(
			make([]zapcore.Field, 0, len(c.fields)+len(fields)),
			c.fields...,
		),
		fields...,
	)

	if e.Level >= c.minLevel {
		enc := zapcore.NewMapObjectEncoder()

		for _, field := range fields {
			field.AddTo(enc)
		}

		attachment := slack.Attachment{
			Color:    levelColor[e.Level],
			Fallback: e.Message,
		}

		attachment.Text = e.Message + "\n"
		for k, v := range enc.Fields {
			attachment.Text += fmt.Sprintf("*%s*\n%v\n", k, v)
		}

		msg := slack.WebhookMessage{
			Attachments: []slack.Attachment{attachment},
		}

		err := slack.PostWebhook(c.hookURL, &msg)

		if err != nil {
			fields = append(fields, zap.Error(err), zap.String("slack_error", "send event to slack error"))
		}
	}

	err := c.Core.Write(e, fields)

	return err
}

func (c *slackCore) With(fields []zapcore.Field) zapcore.Core {
	return &slackCore{
		Core:     c.Core.With(nil),
		hookURL:  c.hookURL,
		minLevel: c.minLevel,
		fields:   append(c.fields, fields...),
	}
}

func (c *slackCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

var levelColor = map[zapcore.Level]string{
	zapcore.DebugLevel: "#9B30FF",
	zapcore.InfoLevel:  "good",
	zapcore.WarnLevel:  "warning",
	zapcore.ErrorLevel: "danger",
	zapcore.FatalLevel: "danger",
	zapcore.PanicLevel: "danger",
}
