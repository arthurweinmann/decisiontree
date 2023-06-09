package decisiontree

import (
	"fmt"

	"github.com/arthurweinmann/decisiontree/internal/settings"
)

type Result struct {
	set  bool
	push bool

	values []any
}

func (q *Result) setSettings(s *settings.Settings) {
}

func (q *Result) run(state *State) (string, error) {
	var vals []any
	for _, vn := range q.values {
		raw, ok := vn.(RawStringValue)
		if ok {
			vals = append(vals, raw)
		} else {
			switch t := vn.(type) {
			default:
				return "", fmt.Errorf("We got an invalid type %T for a value in result", t)
			case string:
				v, ok := state.Values[t]
				if ok {
					vals = append(vals, v)
				}
			}
		}
	}

	if len(vals) == 0 {
		return "", nil
	}

	if q.set {
		if len(vals) == 1 {
			state.Result = vals[0]
		} else {
			state.Result = vals
		}
	} else if q.push {
		switch t := state.Result.(type) {
		default:
			return "", fmt.Errorf("We cannot push into result since result is not an array but a %T", t)
		case []any:
			if len(vals) == 1 {
				state.Result = append(t, vals[0])
			} else {
				state.Result = append(t, vals...)
			}
		}
	}

	return "", nil
}

func SetResult(values ...any) *Result {
	return &Result{
		set:    true,
		values: values,
	}
}

func PushResult(values ...any) *Result {
	return &Result{
		push:   true,
		values: values,
	}
}
