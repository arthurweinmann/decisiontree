package decisiontree

import (
	"fmt"

	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/decisiontree/pkg/options"
)

type ConcatenateNode struct {
	settings *settings.Settings
	opts     []options.Option

	targetName          string
	sep                 string
	valuesToConcatenate []any

	next []DecisionBranch
}

func (q *ConcatenateNode) setSettings(s *settings.Settings) {
	if len(q.opts) > 0 {
		q.settings = s.Clone()
		for _, opt := range q.opts {
			opt.Apply(q.settings)
		}
	} else {
		q.settings = s
	}
}

func (q *ConcatenateNode) run(state *State) ([]DecisionBranch, error) {
	var concatenated string

	for _, vn := range q.valuesToConcatenate {
		var v any

		raw, ok := vn.(RawStringValue)
		if ok {
			v = raw
		} else {
			switch t := vn.(type) {
			default:
				return nil, fmt.Errorf("We got an invalid type %T in concatenate", t)
			case string:
				var ok bool
				v, ok = state.Values[t]
				if !ok {
					return nil, fmt.Errorf("We could not find value %s in State", t)
				}
			}
		}

		switch t := v.(type) {
		default:
			return nil, fmt.Errorf("We only support strings for concatenation and we got %T", t)
		case string:
			if concatenated != "" {
				concatenated += q.sep
			}
			concatenated += t
		}
	}

	state.Values[q.targetName] = concatenated

	return q.next, nil
}

func Concatenate(targetName, sep string, values []any) *ConcatenateNode {
	return &ConcatenateNode{
		sep:                 sep,
		valuesToConcatenate: values,
		targetName:          targetName,
	}
}

func (q *ConcatenateNode) SetOptions(opts ...options.Option) *ConcatenateNode {
	q.opts = opts
	return q
}

func (q *ConcatenateNode) Next(next ...DecisionBranch) *ConcatenateNode {
	q.next = next

	return q
}
