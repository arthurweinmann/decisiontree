package decisiontree

import (
	"fmt"
	"strings"

	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/decisiontree/internal/utils"
	"github.com/arthurweinmann/decisiontree/pkg/options"
	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
)

type QuestionDefinedAnswer struct {
	settings *settings.Settings
	opts     []options.Option

	question        string
	valueName       string
	possibleAnswers []string
	nextPerAnswer   map[string]string
}

func (q *QuestionDefinedAnswer) setSettings(s *settings.Settings) {
	if len(q.opts) > 0 {
		q.settings = s.Clone()
		for _, opt := range q.opts {
			opt.Apply(q.settings)
		}
	} else {
		q.settings = s
	}
}

func (q *QuestionDefinedAnswer) run(state *State) (string, error) {
	if len(q.possibleAnswers) == 0 {
		return "", fmt.Errorf("the question with defined answer %s does not have any associated decision branches", q.question)
	}

	switch {
	default:
		return "", fmt.Errorf("question with defined answer %s does not know which ai api and model to use, please configure it", q.question)
	case q.settings.UseOpenAIModel != "":
		switch q.settings.UseOpenAIModel {
		default:
			return "", fmt.Errorf("for now questions with defined answer do not support openai model %s", q.settings.UseOpenAIModel)
		case openai.GPT3_5_turbo_4k:
			return q.runOpenAIGTP3_5(state)
		}
	}
}

func (q *QuestionDefinedAnswer) runOpenAIGTP3_5(state *State) (string, error) {
	if q.settings.OpenAI == nil || q.settings.OpenAI.APIKey == "" {
		return "", fmt.Errorf("We cannot find an openai api key")
	}

	prompt := `You are a %s machine, you only answer %s. You read a text than a question about this text and you write %s.
	Read between the lines of the text and imagine what a human would say about it.
	Here is the text:
	%s
	Answer by %s to this question about the text:
	%s`

	var machine string
	var withthewords string
	var answerthis string
	var answerby string
	for i := 0; i < len(q.possibleAnswers); i++ {
		if machine != "" {
			machine += "/"
		}
		machine += q.possibleAnswers[i]

		if withthewords != "" {
			withthewords += " or "
		}
		switch q.possibleAnswers[i] {
		default:
			panic("should not happen")
		case "yes", "no":
			withthewords += `with the word "` + q.possibleAnswers[i] + `"`
		case "I don't know":
			withthewords += `with the sentence "` + q.possibleAnswers[i] + `"`
		}

		if answerthis != "" {
			answerthis += " or "
		}
		switch q.possibleAnswers[i] {
		default:
			panic("should not happen")
		case "yes", "no":
			answerthis += `"` + q.possibleAnswers[i] + `" if the answer to the question is ` + q.possibleAnswers[i]
		case "I don't know":
			answerthis += `"I don't know" if you do not know the answer`
		}

		if answerby != "" {
			answerby += " or "
		}
		answerby += q.possibleAnswers[i]
	}

	prompt = fmt.Sprintf(prompt,
		machine,
		withthewords,
		answerthis,
		answerby,
	)

	prompt = utils.FormatPrompt(prompt)

	resp, err := openai.CreateChatCompletion(&openai.ChatCompletionRequest{
		APIKEY: q.settings.OpenAI.APIKey,
		Model:  openai.GPT3_5_turbo_4k,
		Messages: []openai.ChatCompletionMessage{{
			Role:    "user",
			Content: prompt,
		}},
		Temperature: 0.4,
	})
	if err != nil {
		return "", err
	}

	kind, err := isYesOrNoOrIDontKnow(resp.Choices[0].Message.Content, q.possibleAnswers)
	if err != nil {
		resp, err = openai.CreateChatCompletion(&openai.ChatCompletionRequest{
			APIKEY: q.settings.OpenAI.APIKey,
			Model:  openai.GPT3_5_turbo_4k,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
				{
					Role:    "assistant",
					Content: resp.Choices[0].Message.Content,
				},
				{
					Role: "user",
					Content: utils.FormatPrompt(fmt.Sprintf(`We have the following error:
					%s
					Please respond only by %s`, err, answerby)),
				},
			},
			Temperature: 0.4,
		})
		if err != nil {
			return "", err
		}

		kind, err = isYesOrNoOrIDontKnow(resp.Choices[0].Message.Content, q.possibleAnswers)
		if err != nil {
			return "", err
		}
	}

	return q.nextPerAnswer[kind], nil
}

func Question(question, valueName string) *QuestionDefinedAnswer {
	qda := &QuestionDefinedAnswer{
		question:      question,
		valueName:     valueName,
		nextPerAnswer: map[string]string{},
	}

	return qda
}

func (q *QuestionDefinedAnswer) SetOptions(opts ...options.Option) *QuestionDefinedAnswer {
	q.opts = opts
	return q
}

func (q *QuestionDefinedAnswer) Yes(next string) *QuestionDefinedAnswer {
	q.nextPerAnswer["yes"] = next
	q.possibleAnswers = append(q.possibleAnswers, "yes")
	return q
}

func (q *QuestionDefinedAnswer) No(next string) *QuestionDefinedAnswer {
	q.nextPerAnswer["no"] = next
	q.possibleAnswers = append(q.possibleAnswers, "no")
	return q
}

func (q *QuestionDefinedAnswer) DoNotKnow(next string) *QuestionDefinedAnswer {
	q.nextPerAnswer["I don't know"] = next
	q.possibleAnswers = append(q.possibleAnswers, "I don't know")
	return q
}

func isYesOrNoOrIDontKnow(resp string, possiblevalues []string) (string, error) {
	var n, y, idk bool

	for _, pv := range possiblevalues {
		switch pv {
		default:
			panic("should not happen")
		case "yes":
			y = strings.Contains(strings.ToLower(resp), "yes")
		case "no":
			n = strings.Contains(strings.ToLower(resp), "no")
		case "I don't know":
			idk = (strings.Contains(strings.ToLower(resp), "do not") || strings.Contains(strings.ToLower(resp), "don")) && strings.Contains(strings.ToLower(resp), "know")
		}
	}

	var count int
	if n {
		count++
	}
	if y {
		count++
	}
	if idk {
		count++
	}

	if count == 0 {
		return "", fmt.Errorf("response from llm contains neither yes or no or i don't know")
	}
	if count != 1 {
		return "", fmt.Errorf("response from llm contains contradictory answers")
	}

	switch {
	default:
		return "", fmt.Errorf("could not parse response from llm: %s", resp)
	case y:
		return "yes", nil
	case n:
		return "no", nil
	case idk:
		return "I don't know", nil
	}
}
