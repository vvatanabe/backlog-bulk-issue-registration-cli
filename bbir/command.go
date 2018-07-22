package bbir

import (
	"github.com/vvatanabe/go-backlog/backlog/v2"
)

func NewCategoryIDPtr(v int) *CategoryID {
	id := CategoryID(v)
	return &id
}

func NewVersionIDPtr(v int) *VersionID {
	id := VersionID(v)
	return &id
}

func NewUserIDPtr(v int) *UserID {
	id := UserID(v)
	return &id
}

func NewIssueIDPtr(v int) *IssueID {
	id := IssueID(v)
	return &id
}

type ProjectID int
type IssueTypeID int
type PriorityID int
type CategoryID int
type VersionID int
type UserID int
type IssueID int
type CustomFieldID int

type Command struct {
	Line                       int
	ProjectID                  ProjectID
	Summary                    string
	IssueTypeID                IssueTypeID
	PriorityID                 PriorityID
	Description                string
	StartDate                  string
	DueDate                    string
	EstimatedHours             *float64
	ActualHours                *float64
	CategoryID                 *CategoryID
	VersionID                  *VersionID
	MilestoneID                *VersionID
	AssigneeID                 *UserID
	ParentIssueID              *IssueID
	CustomFields               map[CustomFieldID]interface{}
	Children                   []*Command
	HasUnregisteredParentIssue bool
}

func (cmd *Command) AddChild(child *Command) {
	cmd.Children = append(cmd.Children, child)
}

func (cmd *Command) ToAddIssueOptions() *v2.AddIssueOptions {
	opt := &v2.AddIssueOptions{CustomFields: make(map[int]interface{})}
	opt.Description = cmd.Description
	opt.StartDate = cmd.StartDate
	opt.DueDate = cmd.DueDate
	if cmd.EstimatedHours != nil {
		opt.EstimatedHours = *cmd.EstimatedHours
	}
	if cmd.ActualHours != nil {
		opt.ActualHours = *cmd.ActualHours
	}
	if cmd.CategoryID != nil {
		opt.CategoryIDs = []int{int(*cmd.CategoryID)}
	}
	if cmd.VersionID != nil {
		opt.VersionIDs = []int{int(*cmd.VersionID)}
	}
	if cmd.MilestoneID != nil {
		opt.MilestoneIDs = []int{int(*cmd.MilestoneID)}
	}
	if cmd.AssigneeID != nil {
		opt.AssigneeID = int(*cmd.AssigneeID)
	}
	if cmd.ParentIssueID != nil {
		opt.ParentIssueID = int(*cmd.ParentIssueID)
	}
	for k, v := range cmd.CustomFields {
		opt.CustomFields[int(k)] = v
	}
	return opt
}
