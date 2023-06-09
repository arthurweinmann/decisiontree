package decisiontree

import (
	"fmt"

	"github.com/arthurweinmann/decisiontree/internal/settings"
	"github.com/arthurweinmann/decisiontree/pkg/session"
)

const tree_private_internal_key = "_____internal_Do_nOt_Use"

type DecisionBranch interface {
	setSettings(s *settings.Settings)
	run(*State) (string, error)
}

type Tree map[string][]DecisionBranch

type treeAttributes struct {
	lastSession string
}

func (*treeAttributes) setSettings(s *settings.Settings) {}
func (*treeAttributes) run(*State) (string, error)       { return "", nil }

func (t Tree) getAttr() (*treeAttributes, error) {
	attrs, ok := t[tree_private_internal_key]
	if !ok {
		attrs = []DecisionBranch{&treeAttributes{}}
		t[tree_private_internal_key] = attrs
	}
	if len(attrs) != 1 {
		return nil, fmt.Errorf("internal inconsistency")
	}
	attr, ok := attrs[0].(*treeAttributes)
	if !ok {
		return nil, fmt.Errorf("internal inconsistency")
	}
	return attr, nil
}

func (t Tree) initialize(sess *session.Session) error {
	attr, err := t.getAttr()
	if err != nil {
		return err
	}

	if attr.lastSession == "" || attr.lastSession != sess.ID {
		for _, branches := range t {
			for i := 0; i < len(branches); i++ {
				branches[i].setSettings(sess.Settings)
			}
		}

		attr.lastSession = sess.ID
	}

	return nil
}

func (t Tree) RunOnSingleInput(sess *session.Session, valuename string, value any) (any, error) {
	err := t.initialize(sess)
	if err != nil {
		return nil, err
	}

	currentBranchName := "_start"
	currentBranch, ok := t[currentBranchName]
	if !ok {
		return nil, fmt.Errorf("could not find start branch `_start`")
	}

	st := &State{
		Values: map[string]any{
			valuename: value,
		},
	}

	err = t.runBranches(currentBranchName, currentBranch, st)
	if err != nil {
		return nil, err
	}

	return st.Result, err
}

func (t Tree) runBranches(branchName string, branch []DecisionBranch, st *State) error {
	var nextBranches []string
	for j := 0; j < len(branch); j++ {
		next, err := branch[j].run(st)
		if err != nil {
			return fmt.Errorf("We encountered an error in branch %s-%d: %v", branchName, j, err)
		}
		if next != "" {
			nextBranches = append(nextBranches, next)
		}
	}

	// Now recursively process the next branches.
	for _, nextBranchName := range nextBranches {
		nextBranch, ok := t[nextBranchName]
		if !ok {
			return fmt.Errorf("Could not find next branch `%s`", nextBranchName)
		}
		err := t.runBranches(nextBranchName, nextBranch, st)
		if err != nil {
			return err
		}
	}

	return nil
}
