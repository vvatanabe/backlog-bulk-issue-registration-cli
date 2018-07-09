package internal

import (
	"context"

	"sync"

	"github.com/golang/sync/errgroup"
	. "github.com/vvatanabe/go-backlog/backlog/v2"
)

type IssueRepository interface {
	FindIssueByKey(ctx context.Context, issueKey string) (*Issue, error)
	AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error)
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

func (s *IssueHTTPClient) FindIssueByKey(ctx context.Context, issueKey string) (*Issue, error) {
	if v, ok := s.issues.Load(issueKey); ok {
		if v != nil {
			return v.(*Issue), nil
		} else {
			return nil, nil
		}
	}
	if v, err := s.client.GetIssue(ctx, issueKey); err != nil {
		return nil, err
	} else {
		s.issues.Store(issueKey, v)
		return v, nil
	}
}

func (s *IssueHTTPClient) AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error) {
	if v, err := s.client.AddIssue(ctx, projectID, summary, issueTypeID, priorityID, opt); err != nil {
		return nil, err
	} else {
		s.issues.Store(v.IssueKey, v)
		return v, nil
	}
}

type ProjectRepository interface {
	GetProjectID() ProjectID
	FindUserByUserID(userID string) *User
	FindIssueTypeByName(name string) *IssueType
	FindCategoryByName(name string) *Category
	FindVersionByName(name string) *Version
	Prefetch(ctx context.Context) error
}

func NewProjectHTTPClient(cfg *Config, client BacklogAPIClient) ProjectRepository {
	return &ProjectHTTPClient{
		cfg:        cfg,
		client:     client,
		users:      make(map[string]*User),
		issueTypes: make(map[string]*IssueType),
		categories: make(map[string]*Category),
		versions:   make(map[string]*Version),
	}
}

type ProjectHTTPClient struct {
	cfg    *Config
	client BacklogAPIClient

	project    *Project
	users      map[string]*User
	issueTypes map[string]*IssueType
	categories map[string]*Category
	versions   map[string]*Version
}

func (s *ProjectHTTPClient) GetProjectID() ProjectID {
	return ProjectID(s.project.ID)
}

func (s *ProjectHTTPClient) FindUserByUserID(userID string) *User {
	v, ok := s.users[userID]
	if !ok {
		return nil
	}
	return v
}

func (s *ProjectHTTPClient) FindIssueTypeByName(name string) *IssueType {
	v, ok := s.issueTypes[name]
	if !ok {
		return nil
	}
	return v
}

func (s *ProjectHTTPClient) FindCategoryByName(name string) *Category {
	v, ok := s.categories[name]
	if !ok {
		return nil
	}
	return v
}

func (s *ProjectHTTPClient) FindVersionByName(name string) *Version {
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
				s.users[v.UserID] = v
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
