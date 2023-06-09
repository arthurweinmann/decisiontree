package settings

import "github.com/arthurweinmann/go-ai-sdk/pkg/openai"

type Settings struct {
	OpenAI *OpenAISettings

	UseOpenAIModel openai.Model
}

type OpenAISettings struct {
	APIKey string
	Model  openai.Model

	Temperature float32
	MaxTokens   int // -1 for max remaining tokens
}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) Clone() *Settings {
	clone := &Settings{}

	if s.OpenAI != nil {
		clone.OpenAI = &OpenAISettings{
			APIKey:      s.OpenAI.APIKey,
			Model:       s.OpenAI.Model,
			Temperature: s.OpenAI.Temperature,
			MaxTokens:   s.OpenAI.MaxTokens,
		}
	}

	clone.UseOpenAIModel = s.UseOpenAIModel

	return clone
}
