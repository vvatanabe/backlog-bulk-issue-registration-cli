package bbir

import (
	"context"

	. "github.com/vvatanabe/go-backlog/backlog/v2"
)

type BacklogAPIClientAsMock struct {
	addIssue        func(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error)
	getIssue        func(ctx context.Context, issueKey string) (*Issue, error)
	getProject      func(ctx context.Context, projectKey string) (*Project, error)
	getProjectUsers func(ctx context.Context, id ProjectID) ([]*User, error)
	getIssueTypes   func(ctx context.Context, id ProjectID) ([]*IssueType, error)
	getCategories   func(ctx context.Context, id ProjectID) ([]*Category, error)
	getVersions     func(ctx context.Context, id ProjectID) ([]*Version, error)
	getCustomFields func(ctx context.Context, id ProjectID) ([]*CustomField, error)
}

func (m *BacklogAPIClientAsMock) AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error) {
	return m.addIssue(ctx, projectID, summary, issueTypeID, priorityID, opt)
}

func (m *BacklogAPIClientAsMock) GetIssue(ctx context.Context, issueKey string) (*Issue, error) {
	return m.getIssue(ctx, issueKey)
}

func (m *BacklogAPIClientAsMock) GetProject(ctx context.Context, projectKey string) (*Project, error) {
	return m.getProject(ctx, projectKey)
}

func (m *BacklogAPIClientAsMock) GetProjectUsers(ctx context.Context, id ProjectID) ([]*User, error) {
	return m.getProjectUsers(ctx, id)
}

func (m *BacklogAPIClientAsMock) GetIssueTypes(ctx context.Context, id ProjectID) ([]*IssueType, error) {
	return m.getIssueTypes(ctx, id)
}

func (m *BacklogAPIClientAsMock) GetCategories(ctx context.Context, id ProjectID) ([]*Category, error) {
	return m.getCategories(ctx, id)
}

func (m *BacklogAPIClientAsMock) GetVersions(ctx context.Context, id ProjectID) ([]*Version, error) {
	return m.getVersions(ctx, id)
}

func (m *BacklogAPIClientAsMock) GetCustomFields(ctx context.Context, id ProjectID) ([]*CustomField, error) {
	return m.getCustomFields(ctx, id)
}
