package internal

import (
	"context"

	"fmt"

	"github.com/pkg/errors"
	"github.com/vvatanabe/errsgroup"
)

func NewBulkCommandExecutor(issue IssueRepository, msgs Messages) *BulkCommandExecutor {
	return &BulkCommandExecutor{issue, msgs}
}

type BulkCommandExecutor struct {
	issue IssueRepository
	msgs  Messages
}

const Parallelism = 12

func (c BulkCommandExecutor) Do(commands []*Command, callbacks ...Callback) error {

	g := errsgroup.NewGroup(
		errsgroup.MaxParallelSize(Parallelism),
		errsgroup.ErrorChanelSize(Parallelism),
	)

	opts := NewDefaultCallbackOptions()
	for _, cb := range callbacks {
		cb(opts)
	}

	opts.Before()
	for _, v := range commands {
		cmd := v
		g.Go(func() error {
			opt := cmd.ToAddIssueOptions()
			issue, err := c.issue.AddIssue(context.Background(), cmd.ProjectID, cmd.Summary, cmd.IssueTypeID, cmd.PriorityID, opt)
			if err != nil {
				return errors.New(fmt.Sprintf("%v: %v", c.msgs.Line(cmd.Line), err.Error()))
			}
			opts.Each()
			for _, v := range cmd.Children {
				cmd := v
				g.Go(func() error {
					opt := cmd.ToAddIssueOptions()
					opt.ParentIssueID = issue.ID
					_, err := c.issue.AddIssue(context.Background(), cmd.ProjectID, cmd.Summary, cmd.IssueTypeID, cmd.PriorityID, opt)
					if err != nil {
						return errors.New(fmt.Sprintf("%v: %v", c.msgs.Line(cmd.Line), err.Error()))
					}
					opts.Each()
					return nil
				})
			}
			return nil
		})
	}
	if errs := g.Wait(); len(errs) > 0 {
		// TODO sort err
		return NewMultipleErrors(errs)
	}
	opts.After()
	return nil
}
