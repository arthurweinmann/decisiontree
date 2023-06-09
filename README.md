# Decision Tree

An experiment to define and run AI workflow trees, which can potentially be used in a variety of decision-making and logic-routing scenarios.

The workflow decision tree can process user input and decide on a subsequent course of action.

# How to use

```go
package main

import (
    . "github.com/arthurweinmann/decisiontree/pkg/decisiontree"
	"github.com/arthurweinmann/decisiontree/pkg/options"
	"github.com/arthurweinmann/decisiontree/pkg/session"
	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
	"github.com/davecgh/go-spew/spew"
    "log"
)

func main() {
    sess, err := session.NewSession(options.WithOpenAI("YOU_OPENAI_API_KEY"), options.WithOpenAIModel(openai.GPT3_5_turbo_4k), options.WithOpenAIMaxTokens(-1), options.WithOpenAITemperature(0.7))
	if err != nil {
		log.Fatalf("could not create session: %v", err)
	}

	questionAction := Action(func(state *State, act *ActionHandler) (string, error) {
		// For example search in a vector database
		state.Values["prompt"] = "Say Hello World in a cool way"
		return "respond", nil
	})

	mergeLLMRespAction := Action(func(state *State, act *ActionHandler) (string, error) {
		state.Values["resp1"] = fmt.Sprintf("We have the following bug report:\n%s\nThe bug's origin may come from the following:\n%s\nPlease resolve the bug", state.Values["input"], state.Values["resp1"])
		return "setresult", nil
	})

	tree := Tree{
		"_start":       {Question("Is it a question?", "input").Yes("question").No("notaquestion")},
		"question":     {questionAction},
		"notaquestion": {Question("Is it a bug description and/or a debug request?", "input").Yes("debug").No("donotknow")},
		"debug":        {PromptLLM(RawStringValue("imagine possible issues which could be the origin of the bug"), "resp1").Next("mergellmresp")},
		"mergellmresp": {mergeLLMRespAction},
		"respond":      {PromptLLM("prompt", "resp1").SetOptions(options.WithOpenAIMaxTokens(-1), options.WithOpenAITemperature(0.7)).Next("setresult")},
		"donotknow":    {SetResult(RawStringValue("We could not find an answer"))},
		"setresult":    {SetResult("resp1")},
	}

	result, err := tree.RunOnSingleInput(sess, "input", "I have a question")
	if err != nil {
		log.Fatalf("could not run tree: %v", err)
	}

	spew.Dump(result)
}
```