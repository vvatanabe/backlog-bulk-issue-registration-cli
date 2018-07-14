package bbir

import (
	"context"

	"fmt"
	"testing"

	. "github.com/vvatanabe/go-backlog/backlog/v2"
	"github.com/vvatanabe/shot/shot"
)

type IssueRepositoryAsMock struct {
	findIssueByKey func(ctx context.Context, issueKey string) (*Issue, error)
	addIssue       func(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error)
}

func (r *IssueRepositoryAsMock) FindIssueByKey(ctx context.Context, issueKey string) (*Issue, error) {
	return r.findIssueByKey(ctx, issueKey)
}

func (r *IssueRepositoryAsMock) AddIssue(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error) {
	return r.addIssue(ctx, projectID, summary, issueTypeID, priorityID, opt)
}

type ProjectRepositoryAsMock struct {
	getProjectID        func() ProjectID
	findUserByUserID    func(userID string) *User
	findIssueTypeByName func(name string) *IssueType
	findCategoryByName  func(name string) *Category
	findVersionByName   func(name string) *Version
	prefetch            func(ctx context.Context) error
}

func (r *ProjectRepositoryAsMock) GetProjectID() ProjectID {
	return r.getProjectID()
}

func (r *ProjectRepositoryAsMock) FindUserByUserID(userID string) *User {
	return r.findUserByUserID(userID)
}

func (r *ProjectRepositoryAsMock) FindIssueTypeByName(name string) *IssueType {
	return r.findIssueTypeByName(name)
}

func (r *ProjectRepositoryAsMock) FindCategoryByName(name string) *Category {
	return r.findCategoryByName(name)
}

func (r *ProjectRepositoryAsMock) FindVersionByName(name string) *Version {
	return r.findVersionByName(name)
}

func (r *ProjectRepositoryAsMock) Prefetch(ctx context.Context) error {
	return r.prefetch(ctx)
}

func NewInjectorForRepositoryTest(t *testing.T) shot.Injector {
	injector, err := shot.CreateInjector(func(binder shot.Binder) {
		binder.Bind(new(Config)).ToInstance(&Config{ProjectKey: "EXAMPLE"})
		var issueID = 3
		binder.Bind(new(BacklogAPIClient)).ToInstance(&BacklogAPIClientAsMock{
			getProject: func(ctx context.Context, projectKey string) (*Project, error) {
				return &Project{
					ID:         1,
					ProjectKey: projectKey,
				}, nil
			},
			getProjectUsers: func(ctx context.Context, id ProjectID) ([]*User, error) {
				user1 := &User{ID: 1, UserID: "ken"}
				user2 := &User{ID: 2, UserID: "taro"}
				user3 := &User{ID: 3, UserID: "hana"}
				return []*User{user1, user2, user3}, nil
			},
			getIssueTypes: func(ctx context.Context, id ProjectID) ([]*IssueType, error) {
				issueType1 := &IssueType{ID: 1, Name: "task"}
				issueType2 := &IssueType{ID: 2, Name: "bug"}
				issueType3 := &IssueType{ID: 3, Name: "operation"}
				return []*IssueType{issueType1, issueType2, issueType3}, nil
			},
			getCategories: func(ctx context.Context, id ProjectID) ([]*Category, error) {
				category1 := &Category{ID: 1, Name: "web"}
				category2 := &Category{ID: 2, Name: "iPhone"}
				category3 := &Category{ID: 3, Name: "android"}
				return []*Category{category1, category2, category3}, nil
			},
			getVersions: func(ctx context.Context, id ProjectID) ([]*Version, error) {
				version1 := &Version{ID: 1, Name: "sprint1"}
				version2 := &Version{ID: 2, Name: "sprint2"}
				version3 := &Version{ID: 3, Name: "sprint3"}
				return []*Version{version1, version2, version3}, nil
			},
			getIssue: func(ctx context.Context, issueKey string) (*Issue, error) {
				switch issueKey {
				case "EXAMPLE-1":
					return &Issue{ID: 1, IssueKey: issueKey}, nil
				case "EXAMPLE-2":
					return &Issue{ID: 2, IssueKey: issueKey}, nil
				case "EXAMPLE-3":
					return &Issue{ID: 3, IssueKey: issueKey}, nil
				default:
					return nil, nil
				}
			},
			addIssue: func(ctx context.Context, projectID ProjectID, summary string, issueTypeID IssueTypeID, priorityID PriorityID, opt *AddIssueOptions) (*Issue, error) {
				issueID++
				return &Issue{ID: issueID, IssueKey: fmt.Sprintf("EXAMPLE-%v", issueID), ProjectID: int(projectID)}, nil
			},
		})
		binder.Bind(new(ProjectRepository)).ToConstructor(NewProjectHTTPClient).AsEagerSingleton()
		binder.Bind(new(IssueRepository)).ToConstructor(NewIssueHTTPClient).AsEagerSingleton()
	})
	if err != nil {
		t.Errorf("Could not inject test module: %s", err.Error())
	}
	if v, err := injector.SafeGet(new(ProjectRepository)); err != nil {
		t.Errorf("Could not get GetProjectRepository module: %s", err.Error())
	} else {
		if project, ok := v.(ProjectRepository); !ok {
			t.Error("Could not assert type (GetProjectRepository)")
		} else {
			if err := project.Prefetch(context.Background()); err != nil {
				t.Errorf("Could not prefetch: %s", err.Error())
			}
		}
	}
	return injector
}

func Test_IssueHTTPClient_FindIssueByKey_should_return_issue_that_match_name(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	issue := injector.Get(new(IssueRepository)).(*IssueHTTPClient)
	want := "EXAMPLE-1" // TODO change table driven test
	result, err := issue.FindIssueByKey(context.Background(), want)
	if err != nil {
		t.Errorf("Could not match result: %s", err.Error())
	}
	if result == nil {
		t.Error("Result is nil")
	}
	if want != result.IssueKey {
		t.Errorf("Could not match result. want: %v, result:  %v", want, result)
	}
}

func Test_ProjectHTTPClient_GetProjectID(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)

	want := ProjectID(1)
	result := project.GetProjectID()
	if want != result {
		t.Errorf("Could not match result. want: %v, result:  %v", want, result)
	}
}

func Test_ProjectHTTPClient_FindUserByUserID_should_return_user_that_match_name(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	want := "ken"
	result := project.FindUserByUserID(want)
	if result == nil {
		t.Error("Result is nil")
	}
	if want != result.UserID {
		t.Errorf("Could not match result. want: %v, result:  %v", want, result)
	}
}

func Test_ProjectHTTPClient_FindUserByUserID_should_return_nil_if_does_not_match_userID(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	want := "xxx"
	result := project.FindUserByUserID(want)
	if result != nil {
		t.Error("Result is not nil")
	}
}

func Test_ProjectHTTPClient_FindIssueTypeByName_should_return_issueType_that_match_name(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	want := "task"
	result := project.FindIssueTypeByName(want)
	if result == nil {
		t.Error("Result is nil")
	}
	if want != result.Name {
		t.Errorf("Could not match result. want: %v, result:  %v", want, result)
	}
}

func Test_ProjectHTTPClient_FindIssueTypeByName_should_return_nil_if_does_not_match_ID(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	result := project.FindIssueTypeByName("xxx")
	if result != nil {
		t.Error("Result is not nil")
	}
}

func Test_ProjectHTTPClient_FindCategoryByName_should_return_category_that_match_name(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	want := "web"
	result := project.FindCategoryByName(want)
	if result == nil {
		t.Error("Result is nil")
	}
	if want != result.Name {
		t.Errorf("Could not match result. want: %v, result:  %v", want, result)
	}
}

func Test_ProjectHTTPClient_FindCategoryByName_should_return_nil_if_does_not_match_name(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	result := project.FindCategoryByName("xxx")
	if result != nil {
		t.Error("Result is not nil")
	}
}

func Test_ProjectHTTPClient_FindVersionByName_should_return_version_that_match_name(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	want := "sprint1"
	result := project.FindVersionByName(want)
	if result == nil {
		t.Error("Result is nil")
	}
	if want != result.Name {
		t.Errorf("Could not match result. want: %v, result:  %v", want, result)
	}
}

func Test_ProjectHTTPClient_FindVersionByName_should_return_nil_if_does_not_match_name(t *testing.T) {
	injector := NewInjectorForRepositoryTest(t)
	project := injector.Get(new(ProjectRepository)).(*ProjectHTTPClient)
	result := project.FindVersionByName("xxx")
	if result != nil {
		t.Error("Result is not nil")
	}
}
