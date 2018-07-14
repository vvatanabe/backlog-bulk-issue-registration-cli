package bbir

import (
	"context"
	"fmt"

	. "github.com/vvatanabe/go-backlog/backlog/v2"
)

type BacklogAPIClient interface {
	AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error)
	GetIssue(ctx context.Context, issueKey string) (*Issue, error)
	GetProject(ctx context.Context, projectKey string) (*Project, error)
	GetProjectUsers(ctx context.Context, id ProjectID) ([]*User, error)
	GetIssueTypes(ctx context.Context, id ProjectID) ([]*IssueType, error)
	GetCategories(ctx context.Context, id ProjectID) ([]*Category, error)
	GetVersions(ctx context.Context, id ProjectID) ([]*Version, error)
}

func NewBacklogAPIClient(cfg *Config) BacklogAPIClient {
	c := NewClient(cfg.SpaceDomain, nil)
	c.SetAPIKey(cfg.APIKey)
	return &client{c}
}

type client struct {
	*Client
}

func (c *client) AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error) {
	v, _, err := c.Issues.AddIssue(ctx, int(projectID), summary, int(issueTypeID), int(priorityID), opt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetIssue(ctx context.Context, issueKey string) (*Issue, error) {
	v, _, err := c.Issues.GetIssue(ctx, issueKey)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetProject(ctx context.Context, projectKey string) (*Project, error) {
	v, _, err := c.Projects.GetProject(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetProjectUsers(ctx context.Context, id ProjectID) ([]*User, error) {
	v, _, err := c.Projects.GetProjectUsers(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetIssueTypes(ctx context.Context, id ProjectID) ([]*IssueType, error) {
	v, _, err := c.Projects.GetIssueTypes(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetCategories(ctx context.Context, id ProjectID) ([]*Category, error) {
	v, _, err := c.Projects.GetCategories(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetVersions(ctx context.Context, id ProjectID) ([]*Version, error) {
	v, _, err := c.Projects.GetVersions(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}
