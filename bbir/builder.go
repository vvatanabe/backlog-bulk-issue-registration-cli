package bbir

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

type CommandBuilder interface {
	Build(int, *Line) (*Command, error)
}

func NewCommandBuilder(issue IssueRepository, project ProjectRepository, msgs Messages) CommandBuilder {
	return &commandBuilder{issue, project, msgs}
}

type commandBuilder struct {
	issue   IssueRepository
	project ProjectRepository
	msgs    Messages
}

type resolver func(*Command, *Line) error

func (b *commandBuilder) Build(lineNum int, line *Line) (*Command, error) {
	var errs []error
	command := &Command{Line: lineNum, Children: []*Command{}}
	for _, f := range []resolver{
		b.ensureProjectID,
		b.ensureSummary,
		b.ensureIssueTypeID,
		b.ensurePriorityID,
		b.ensureDescription,
		b.ensureStartAndDueDate,
		b.ensureEstimatedHours,
		b.ensureActualHours,
		b.ensureCategoryID,
		b.ensureVersionID,
		b.ensureMilestoneID,
		b.ensureAssigneeID,
		b.ensureParentIssue,
	} {
		if err := f(command, line); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return nil, NewMultipleErrors(errs)
	}
	return command, nil
}

func (b *commandBuilder) ensureProjectID(cmd *Command, line *Line) error {
	cmd.ProjectID = b.project.GetProjectID()
	return nil
}

func (b *commandBuilder) ensureSummary(cmd *Command, line *Line) error {
	if line.Summary == "" {
		return errors.New(b.msgs.SummaryIsRequired())
	}
	cmd.Summary = line.Summary
	return nil
}

func (b *commandBuilder) ensureDescription(cmd *Command, line *Line) error {
	cmd.Description = line.Description
	return nil
}

const DateLayout = "2006-01-02"

func (b *commandBuilder) ensureStartAndDueDate(cmd *Command, line *Line) error {
	var startDate, dueDate time.Time
	if line.StartDate != "" {
		var err error
		startDate, err = time.Parse(DateLayout, line.StartDate)
		if err != nil {
			return errors.New(b.msgs.StartDateIsInvalid(line.StartDate))
		}
	}
	if line.DueDate != "" {
		var err error
		dueDate, err = time.Parse(DateLayout, line.DueDate)
		if err != nil {
			return errors.New(b.msgs.DueDateIsInvalid(line.DueDate))
		}
	}
	if !startDate.IsZero() && !dueDate.IsZero() {
		if startDate.After(dueDate) {
			return errors.New(b.msgs.StartDateIsAfterDueDate(line.StartDate, line.DueDate))
		}
	}
	cmd.StartDate = line.StartDate
	cmd.DueDate = line.DueDate
	return nil
}

func (b *commandBuilder) ensureEstimatedHours(cmd *Command, line *Line) error {
	if line.EstimatedHours != "" {
		split := strings.Split(line.EstimatedHours, ".")
		// Pattern Ex. "1.2.3" => ["1", "2", "3"]
		if len(split) > 2 {
			return errors.New(b.msgs.EstimatedHoursIsInvalid(line.EstimatedHours))
		}
		// Pattern Ex. 1.234 => ["1", "234"]
		if len(split) == 2 && len(split[1]) > 2 {
			return errors.New(b.msgs.EstimatedHoursIsInvalid(line.EstimatedHours))
		}
		v, err := strconv.ParseFloat(line.EstimatedHours, 64)
		if err != nil || (math.Signbit(v) && v < 0) {
			return errors.New(b.msgs.EstimatedHoursIsInvalid(line.EstimatedHours))
		}
		cmd.EstimatedHours = &v
	}
	return nil
}

func (b *commandBuilder) ensureActualHours(cmd *Command, line *Line) error {
	if line.ActualHours != "" {
		split := strings.Split(line.ActualHours, ".")
		// Pattern Ex. "1.2.3" => ["1", "2", "3"]
		if len(split) > 2 {
			return errors.New(b.msgs.EstimatedHoursIsInvalid(line.ActualHours))
		}
		// Pattern Ex. 1.234 => ["1", "234"]
		if len(split) == 2 && len(split[1]) > 2 {
			return errors.New(b.msgs.EstimatedHoursIsInvalid(line.ActualHours))
		}
		v, err := strconv.ParseFloat(line.ActualHours, 64)
		if err != nil || (math.Signbit(v) && v < 0) {
			return errors.New(b.msgs.EstimatedHoursIsInvalid(line.ActualHours))
		}
		cmd.ActualHours = &v
	}
	return nil
}

func (b *commandBuilder) ensureIssueTypeID(cmd *Command, line *Line) error {
	issueType := b.project.FindIssueTypeByName(line.IssueType)
	if issueType == nil {
		return errors.New(b.msgs.IssueTypeIsRequired())
	}
	cmd.IssueTypeID = IssueTypeID(issueType.ID)
	return nil
}

const DefaultPriorityID = 3

func (b *commandBuilder) ensurePriorityID(cmd *Command, line *Line) error {
	priorityID := DefaultPriorityID
	if line.Priority != "" {
		v, _ := strconv.Atoi(line.Priority)
		priorityID = v
	}
	// TODO Get PriorityID via project repository
	if !(2 <= priorityID && priorityID <= 4) {
		return errors.New(b.msgs.PriorityIsInvalid(line.Priority))
	}
	cmd.PriorityID = PriorityID(priorityID)
	return nil
}

func (b *commandBuilder) ensureCategoryID(cmd *Command, line *Line) error {
	if line.Category != "" {
		if v := b.project.FindCategoryByName(line.Category); v != nil {
			cmd.CategoryID = NewCategoryIDPtr(v.ID)
		} else {
			return errors.New(b.msgs.CategoryIsNotRegistered(line.Category))
		}
	}
	return nil
}

func (b *commandBuilder) ensureVersionID(cmd *Command, line *Line) error {
	if line.Version != "" {
		v := b.project.FindVersionByName(line.Version)
		if v == nil {
			return errors.New(b.msgs.VersionIsNotRegistered(line.Version))
		}
		cmd.VersionID = NewVersionIDPtr(v.ID)
	}
	return nil
}

func (b *commandBuilder) ensureMilestoneID(cmd *Command, line *Line) error {
	if line.Milestone != "" {
		v := b.project.FindVersionByName(line.Milestone)
		if v == nil {
			return errors.New(b.msgs.MilestoneIsNotRegistered(line.Milestone))
		}
		cmd.MilestoneID = NewVersionIDPtr(v.ID)
	}
	return nil
}

func (b *commandBuilder) ensureAssigneeID(cmd *Command, line *Line) error {
	if line.Assignee != "" {
		v := b.project.FindUserByUserID(line.Assignee)
		if v == nil {
			return errors.New(b.msgs.AssigneeIsNotJoining(line.Assignee))
		}
		cmd.AssigneeID = NewUserIDPtr(v.ID)
	}
	return nil
}

func (b *commandBuilder) ensureParentIssue(cmd *Command, line *Line) error {
	if line.ParentIssue == "" {
		return nil
	}
	if line.ParentIssue == "*" {
		cmd.HasUnregisteredParentIssue = true
		return nil
	}
	if v, err := b.issue.FindIssueByKey(context.Background(), line.ParentIssue); err != nil || v == nil {
		return errors.New(b.msgs.ParentIssueIsNotRegistered(line.ParentIssue))
	} else if v.ParentIssueID != nil {
		return errors.New(b.msgs.ParentIssueAlreadyRegisteredAsChildIssue(line.ParentIssue))
	} else {
		cmd.ParentIssueID = NewIssueIDPtr(v.ID)
		return nil
	}
}
