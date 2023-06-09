package session

import (
	"fmt"

	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/decisiontree/pkg/options"
	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
	"github.com/rs/xid"
)

type Session struct {
	ID       string
	Settings *settings.Settings
}

func NewSession(opts ...options.Option) (*Session, error) {
	sess := &Session{
		ID:       xid.New().String(),
		Settings: settings.NewSettings(),
	}

	for _, opt := range opts {
		opt.Apply(sess.Settings)
	}

	if sess.Settings.OpenAI != nil {
		if sess.Settings.OpenAI.APIKey == "" {
			return nil, fmt.Errorf("openai api key is empty")
		}

		err := openai.Init(sess.Settings.OpenAI.APIKey)
		if err != nil {
			return nil, fmt.Errorf("We could not setup openai sdk: %v", err)
		}
	}

	return sess, nil
}
