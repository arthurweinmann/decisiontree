package decisiontree

import (
	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/decisiontree/pkg/options"
)

type ActionHandlerFunc func(*State, *ActionHandler) (string, error)

type ActionHandler struct {
	settings *settings.Settings
	opts     []options.Option
	fn       ActionHandlerFunc
}

func (q *ActionHandler) setSettings(s *settings.Settings) {
	if len(q.opts) > 0 {
		q.settings = s.Clone()
		for _, opt := range q.opts {
			opt.Apply(q.settings)
		}
	} else {
		q.settings = s
	}
}

func (q *ActionHandler) run(state *State) (string, error) {
	nextbranch, err := q.fn(state, q)
	if err != nil {
		return "", err
	}

	return nextbranch, nil
}

func Action(handler ActionHandlerFunc) *ActionHandler {
	return &ActionHandler{
		fn: handler,
	}
}
