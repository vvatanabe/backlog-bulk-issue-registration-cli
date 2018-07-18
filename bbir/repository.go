package bbir

import (
	"context"

	"sync"

	"github.com/golang/sync/errgroup"
	"github.com/vvatanabe/go-backlog/backlog/v2"
)

type IssueRepository interface {
	FindIssueByKey(ctx context.Context, issueKey string) (*v2.Issue, error)
	AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *v2.AddIssueOptions) (*v2.Issue, error)
}

func NewIssueHTTPClient(client BacklogAPIClient) IssueRepository {
	return &IssueHTTPClient{
		client: client,
		issues: sync.Map{},
	}
}

type IssueHTTPClient struct {
	client BacklogAPIClient
	issues sync.Map
}

func (s *IssueHTTPClient) FindIssueByKey(ctx context.Context, issueKey string) (*v2.Issue, error) {
	if v, ok := s.issues.Load(issueKey); ok {
		if v != nil {
			return v.(*v2.Issue), nil
		}
		return nil, nil
	}
	v, err := s.client.GetIssue(ctx, issueKey)
	if err != nil {
		return nil, err
	}
	s.issues.Store(issueKey, v)
	return v, nil
}

func (s *IssueHTTPClient) AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *v2.AddIssueOptions) (*v2.Issue, error) {
	v, err := s.client.AddIssue(ctx, projectID, summary, issueTypeID, priorityID, opt)
	if err != nil {
		return nil, err
	}
	s.issues.Store(v.IssueKey, v)
	return v, nil
}

type ProjectRepository interface {
	GetProjectID() ProjectID
	FindUserByName(name string) *v2.User
	FindIssueTypeByName(name string) *v2.IssueType
	FindCategoryByName(name string) *v2.Category
	FindVersionByName(name string) *v2.Version
	Prefetch(ctx context.Context) error
}

func NewProjectHTTPClient(cfg *Config, client BacklogAPIClient) ProjectRepository {
	return &ProjectHTTPClient{
		cfg:        cfg,
		client:     client,
		users:      make(map[string]*v2.User),
		issueTypes: make(map[string]*v2.IssueType),
		categories: make(map[string]*v2.Category),
		versions:   make(map[string]*v2.Version),
	}
}

type ProjectHTTPClient struct {
	cfg    *Config
	client BacklogAPIClient

	project    *v2.Project
	users      map[string]*v2.User
	issueTypes map[string]*v2.IssueType
	categories map[string]*v2.Category
	versions   map[string]*v2.Version
}

func (s *ProjectHTTPClient) GetProjectID() ProjectID {
	return ProjectID(s.project.ID)
}

func (s *ProjectHTTPClient) FindUserByName(name string) *v2.User {
	v, ok := s.users[name]
	if !ok {
		return nil
	}
	return v
}

func (s *ProjectHTTPClient) FindIssueTypeByName(name string) *v2.IssueType {
	v, ok := s.issueTypes[name]
	if !ok {
		return nil
	}
	return v
}

func (s *ProjectHTTPClient) FindCategoryByName(name string) *v2.Category {
	v, ok := s.categories[name]
	if !ok {
		return nil
	}
	return v
}

func (s *ProjectHTTPClient) FindVersionByName(name string) *v2.Version {
	v, ok := s.versions[name]
	if !ok {
		return nil
	}
	return v
}

func (s *ProjectHTTPClient) Prefetch(ctx context.Context) error {
	project, err := s.client.GetProject(ctx, s.cfg.ProjectKey)
	if err != nil {
		return err
	}
	s.project = project
	id := ProjectID(s.project.ID)

	g, errCtx := errgroup.WithContext(ctx)
	for _, f := range []func() error{
		func() error {
			users, err := s.client.GetProjectUsers(errCtx, id)
			if err != nil {
				return err
			}
			for _, v := range users {
				s.users[v.Name] = v
			}
			return nil
		},
		func() error {
			issueTypes, err := s.client.GetIssueTypes(errCtx, id)
			if err != nil {
				return err
			}
			for _, v := range issueTypes {
				s.issueTypes[v.Name] = v
			}
			return nil
		},
		func() error {
			categories, err := s.client.GetCategories(errCtx, id)
			if err != nil {
				return err
			}
			for _, v := range categories {
				s.categories[v.Name] = v
			}
			return nil
		},
		func() error {
			versions, err := s.client.GetVersions(errCtx, id)
			if err != nil {
				return err
			}
			for _, v := range versions {
				s.versions[v.Name] = v
			}
			return nil
		},
	} {
		g.Go(f)
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
