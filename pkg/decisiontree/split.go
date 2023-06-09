package decisiontree

import (
	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/decisiontree/pkg/options"
)

type Split struct {
	settings *settings.Settings
	opts     []options.Option

	splitBy         string // ",", "whitespace", ...
	value           any
	targetValueName string

	next string
}

func (q *Split) setSettings(s *settings.Settings) {
	if len(q.opts) > 0 {
		q.settings = s.Clone()
		for _, opt := range q.opts {
			opt.Apply(q.settings)
		}
	} else {
		q.settings = s
	}
}

func (q *Split) run(state *State) (string, error) {
	return q.next, nil
}

func SplitBy(splitby string, valuetosplit any, targetValueName string) *Split {
	return &Split{
		splitBy:         splitby,
		value:           valuetosplit,
		targetValueName: targetValueName,
	}
}

func (q *Split) SetOptions(opts ...options.Option) *Split {
	q.opts = opts
	return q
}

func (q *Split) Next(next string) *Split {
	q.next = next
	return q
}
