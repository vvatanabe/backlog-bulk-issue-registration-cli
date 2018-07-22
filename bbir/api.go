package bbir

import (
	"context"
	"fmt"

	"github.com/vvatanabe/go-backlog/backlog/v2"
)

type BacklogAPIClient interface {
	AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *v2.AddIssueOptions) (*v2.Issue, error)
	GetIssue(ctx context.Context, issueKey string) (*v2.Issue, error)
	GetProject(ctx context.Context, projectKey string) (*v2.Project, error)
	GetProjectUsers(ctx context.Context, id ProjectID) ([]*v2.User, error)
	GetIssueTypes(ctx context.Context, id ProjectID) ([]*v2.IssueType, error)
	GetCategories(ctx context.Context, id ProjectID) ([]*v2.Category, error)
	GetVersions(ctx context.Context, id ProjectID) ([]*v2.Version, error)
	GetCustomFields(ctx context.Context, id ProjectID) ([]*v2.CustomField, error)
}

func NewBacklogAPIClient(cfg *Config) BacklogAPIClient {
	c := v2.NewClient(cfg.SpaceDomain, nil)
	c.SetAPIKey(cfg.APIKey)
	return &client{c}
}

type client struct {
	*v2.Client
}

func (c *client) AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *v2.AddIssueOptions) (*v2.Issue, error) {
	v, _, err := c.Issues.AddIssue(ctx, int(projectID), summary, int(issueTypeID), int(priorityID), opt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetIssue(ctx context.Context, issueKey string) (*v2.Issue, error) {
	v, _, err := c.Issues.GetIssue(ctx, issueKey)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetProject(ctx context.Context, projectKey string) (*v2.Project, error) {
	v, _, err := c.Projects.GetProject(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetProjectUsers(ctx context.Context, id ProjectID) ([]*v2.User, error) {
	v, _, err := c.Projects.GetProjectUsers(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetIssueTypes(ctx context.Context, id ProjectID) ([]*v2.IssueType, error) {
	v, _, err := c.Projects.GetIssueTypes(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetCategories(ctx context.Context, id ProjectID) ([]*v2.Category, error) {
	v, _, err := c.Projects.GetCategories(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetVersions(ctx context.Context, id ProjectID) ([]*v2.Version, error) {
	v, _, err := c.Projects.GetVersions(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *client) GetCustomFields(ctx context.Context, id ProjectID) ([]*v2.CustomField, error) {
	v, _, err := c.Projects.GetCustomFields(ctx, fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	return v, nil
}
