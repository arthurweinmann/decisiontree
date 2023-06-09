package decisiontree

import (
	"fmt"

	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/decisiontree/pkg/options"
	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
)

type PromptChatLLM struct {
	settings *settings.Settings
	opts     []options.Option

	promptValue   any
	respValueName string

	next string
}

func (q *PromptChatLLM) setSettings(s *settings.Settings) {
	if len(q.opts) > 0 {
		q.settings = s.Clone()
		for _, opt := range q.opts {
			opt.Apply(q.settings)
		}
	} else {
		q.settings = s
	}
}

func (q *PromptChatLLM) run(state *State) (string, error) {
	var vs string
	rvs, ok := q.promptValue.(RawStringValue)
	if ok {
		vs = string(rvs)
	} else {
		switch t := q.promptValue.(type) {
		default:
			return "", fmt.Errorf("We got an invalid type %T instead of a string as a value name", t)
		case string:
			v, ok := state.Values[t]
			if !ok {
				return "", fmt.Errorf("We could not find %s value in state", t)
			}
			vs, ok = v.(string)
			if !ok {
				return "", fmt.Errorf("Prompt value is not a string but a %T", v)
			}
		}
	}

	switch {
	default:
		return "", fmt.Errorf("Please configure which ai api and model to use")
	case q.settings.UseOpenAIModel != "":
		switch q.settings.UseOpenAIModel {
		default:
			return "", fmt.Errorf("for now questions with defined answer do not support openai model %s", q.settings.UseOpenAIModel)
		case openai.GPT3_5_turbo_4k:
			return q.runOpenAIGTP3_5(vs, state)
		}
	}
}

func (q *PromptChatLLM) runOpenAIGTP3_5(prompt string, state *State) (string, error) {
	var err error

	if q.settings.OpenAI == nil || q.settings.OpenAI.APIKey == "" {
		return "", fmt.Errorf("We cannot find an openai api key")
	}

	maxtokens := q.settings.OpenAI.MaxTokens
	if maxtokens <= 0 {
		maxtokens, err = openai.GetMaxRemainingTokens(prompt, q.settings.UseOpenAIModel)
		if err != nil {
			return "", err
		}
	}

	resp, err := openai.CreateChatCompletion(&openai.ChatCompletionRequest{
		APIKEY:      q.settings.OpenAI.APIKey,
		Model:       q.settings.UseOpenAIModel,
		MaxTokens:   maxtokens,
		Temperature: float32(q.settings.OpenAI.Temperature),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	state.Values[q.respValueName] = resp.Choices[0].Message.Content

	return q.next, nil
}

func PromptLLM(promptValue any, respValueName string) *PromptChatLLM {
	return &PromptChatLLM{
		promptValue:   promptValue,
		respValueName: respValueName,
	}
}

func (q *PromptChatLLM) SetOptions(opts ...options.Option) *PromptChatLLM {
	q.opts = opts
	return q
}

func (q *PromptChatLLM) Next(next string) *PromptChatLLM {
	q.next = next

	return q
}
