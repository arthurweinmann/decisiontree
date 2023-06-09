package options

import (
	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
)

type Option interface {
	Apply(*settings.Settings)
}

type withOpenAI string

func (w withOpenAI) Apply(o *settings.Settings) {
	if o.OpenAI == nil {
		o.OpenAI = &settings.OpenAISettings{}
	}
	o.OpenAI.APIKey = string(w)
}

func WithOpenAI(apikey string) Option {
	return withOpenAI(apikey)
}

type withOpenAIModel openai.Model

func (w withOpenAIModel) Apply(o *settings.Settings) {
	// TODO: when there are other providers, empty them all here except openai
	o.UseOpenAIModel = openai.Model(w)
}

func WithOpenAIModel(model openai.Model) Option {
	return withOpenAIModel(model)
}

type withOpenAITemperature float32

func (w withOpenAITemperature) Apply(o *settings.Settings) {
	if o.OpenAI == nil {
		o.OpenAI = &settings.OpenAISettings{}
	}
	o.OpenAI.Temperature = float32(w)
}

func WithOpenAITemperature(temperature float32) Option {
	return withOpenAITemperature(temperature)
}

type withOpenAIMaxTokens int

func (w withOpenAIMaxTokens) Apply(o *settings.Settings) {
	if o.OpenAI == nil {
		o.OpenAI = &settings.OpenAISettings{}
	}
	o.OpenAI.MaxTokens = int(w)
}

func WithOpenAIMaxTokens(maxtokens int) Option {
	return withOpenAIMaxTokens(maxtokens)
}
