package bbir

import (
	"context"
	"testing"

	"fmt"

	"sync/atomic"

	. "github.com/vvatanabe/go-backlog/backlog/v2"
	"github.com/vvatanabe/shot/shot"
)

func NewInjectorForCommandBuilderTest(t *testing.T) shot.Injector {
	injector, err := shot.CreateInjector(func(binder shot.Binder) {
		binder.Bind(new(Config)).ToInstance(&Config{
			ProjectKey: "EXAMPLE",
		})
		binder.Bind(new(BulkCommandExecutor)).ToConstructor(NewBulkCommandExecutor).AsEagerSingleton()
		binder.Bind(new(CommandConverter)).ToConstructor(NewCommandConverter).AsEagerSingleton()
		binder.Bind(new(CommandBuilder)).ToConstructor(NewCommandBuilder).AsEagerSingleton()
		binder.Bind(new(Messages)).ToConstructor(NewJapanese).AsEagerSingleton()
		var issueID uint64 = 3
		binder.Bind(new(BacklogAPIClient)).ToInstance(&BacklogAPIClientAsMock{
			getProject: func(ctx context.Context, projectKey string) (*Project, error) {
				return &Project{
					ID:         1,
					ProjectKey: projectKey,
				}, nil
			},
			getProjectUsers: func(ctx context.Context, id ProjectID) ([]*User, error) {
				user1 := &User{ID: 1, Name: "ken"}
				user2 := &User{ID: 2, Name: "taro"}
				user3 := &User{ID: 3, Name: "hana"}
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
			getCustomFields: func(ctx context.Context, id ProjectID) ([]*CustomField, error) {
				text := &CustomField{ID: 1, TypeID: 1, Name: "Text"}
				sentence := &CustomField{ID: 2, TypeID: 2, Name: "Sentence"}
				number := &CustomField{ID: 3, TypeID: 3, Name: "Number"}
				date := &CustomField{ID: 4, TypeID: 4, Name: "Date"}
				singleList := &CustomField{ID: 5, TypeID: 5, Name: "Single List", Items: []*CustomFieldItem{
					{ID: 1, Name: "select-1"},
					{ID: 2, Name: "select-2"},
					{ID: 3, Name: "select-3"},
				}}
				multipleList := &CustomField{ID: 6, TypeID: 6, Name: "Multiple List", Items: []*CustomFieldItem{
					{ID: 1, Name: "select-1"},
					{ID: 2, Name: "select-2"},
					{ID: 3, Name: "select-3"},
				}}
				checkbox := &CustomField{ID: 7, TypeID: 7, Name: "Checkbox", Items: []*CustomFieldItem{
					{ID: 1, Name: "select-1"},
					{ID: 2, Name: "select-2"},
					{ID: 3, Name: "select-3"},
				}}
				radio := &CustomField{ID: 8, TypeID: 8, Name: "Radio", Items: []*CustomFieldItem{
					{ID: 1, Name: "select-1"},
					{ID: 2, Name: "select-2"},
					{ID: 3, Name: "select-3"},
				}}
				return []*CustomField{text, sentence, number, date, singleList, multipleList, checkbox, radio}, nil
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
				issueID := int(atomic.AddUint64(&issueID, 1))
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

func Test_CommandBuilder(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	if _, err := builder.Build(1, &Line{
		Summary:        "Summary",
		Description:    "Description",
		StartDate:      "2018-01-01",
		DueDate:        "2018-01-02",
		EstimatedHours: "1",
		ActualHours:    "2",
		IssueType:      "task",
		Category:       "web",
		Version:        "sprint1",
		Milestone:      "sprint1",
		Priority:       "2",
		Assignee:       "ken",
		ParentIssue:    "*",
	}); err != nil {
		t.Errorf("Could not build a command: %s", err.Error())
	}
}

func Test_CommandBuilder_resolveSummary(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	command := &Command{}
	line := &Line{
		Summary: "Summary",
	}
	if err := builder.ensureSummary(command, line); err != nil {
		t.Errorf("Could not resolve. %s", err.Error())
	} else if command.Summary != line.Summary {
		t.Errorf("Could not resolve a value. want: %s, result: %s", line.Summary, command.Summary)
	}
}

func Test_CommandBuilder_resolveSummary_should_return_error_if_summary_is_empty(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	command := &Command{}
	line := &Line{}
	if err := builder.ensureSummary(command, line); err == nil {
		t.Errorf("Should return error if summary is empty")
	}
}

func Test_CommandBuilder_resolveDescription(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	command := &Command{}
	line := &Line{
		Description: "Description",
	}
	if err := builder.ensureDescription(command, line); err != nil {
		t.Errorf("Could not resolve. %s", err.Error())
	} else if command.Description != line.Description {
		t.Errorf("Could not resolve a value. want: %s, result: %s", line.Description, command.Description)
	}
}

func Test_CommandBuilder_resolveStartAndDueDate(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []struct {
		startDate string
		dueDate   string
	}{
		{
			startDate: "2018-01-01",
			dueDate:   "2018-01-02",
		},
		{
			startDate: "2018-11-01",
			dueDate:   "2018-12-01",
		},
		{
			startDate: "2018-11-11",
			dueDate:   "2018-12-12",
		},
		{
			startDate: "2018-11-01",
			dueDate:   "2018-11-01",
		},
		{
			startDate: "",
			dueDate:   "",
		},
		{
			startDate: "2018-11-01",
			dueDate:   "",
		},
		{
			startDate: "",
			dueDate:   "2018-11-01",
		},
	}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			StartDate: v.startDate,
			DueDate:   v.dueDate,
		}
		if err := builder.ensureStartAndDueDate(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		} else if command.Description != line.Description {
			t.Errorf("Could not resolve a value. want: %s, result: %s", line.Description, command.Description)
		}
	}
}

func Test_CommandBuilder_resolveStartAndDueDate_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []struct {
		startDate string
		dueDate   string
	}{
		{
			startDate: "2018-11-01",
			dueDate:   "2018-10-01",
		},
		{
			startDate: "2018-11-01",
			dueDate:   "2017-10-01",
		},
	}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			StartDate: v.startDate,
			DueDate:   v.dueDate,
		}
		if err := builder.ensureSummary(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolveEstimatedHours(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "1", "1.0", "1.25", "0.50"}
	for _, v := range tests {

		command := &Command{}
		line := &Line{
			EstimatedHours: v,
		}
		if err := builder.ensureEstimatedHours(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveEstimatedHours_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"abc", "1a", "1.a", "1.5a", "1.1.1", "1.111"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			EstimatedHours: v,
		}
		if err := builder.ensureEstimatedHours(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolveActualHours(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "1", "1.0", "1.25", "0.50"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			ActualHours: v,
		}
		if err := builder.ensureActualHours(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveActualHours_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"abc", "1a", "1.a", "1.5a", "1.1.1", "1.111"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			ActualHours: v,
		}
		if err := builder.ensureActualHours(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolveIssueTypeID(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"task", "bug"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			IssueType: v,
		}
		if err := builder.ensureIssueTypeID(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveIssueTypeID_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "xxx"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			IssueType: v,
		}
		if err := builder.ensureIssueTypeID(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolveCategoryID(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "web", "iPhone"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Category: v,
		}
		if err := builder.ensureCategoryID(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveCategoryID_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"xxx"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Category: v,
		}
		if err := builder.ensureCategoryID(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolveVersionID(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "sprint1"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Version: v,
		}
		if err := builder.ensureVersionID(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveVersionID_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"xxx"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Version: v,
		}
		if err := builder.ensureVersionID(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolveMilestoneID(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "sprint1"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Milestone: v,
		}
		if err := builder.ensureMilestoneID(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveMilestoneID_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"xxx"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Milestone: v,
		}
		if err := builder.ensureMilestoneID(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolveAssigneeID(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "ken"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Assignee: v,
		}
		if err := builder.ensureAssigneeID(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveAssigneeID_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"xxx"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Assignee: v,
		}
		if err := builder.ensureAssigneeID(command, line); err == nil {
			t.Errorf("Should return error.")
		}
	}
}

func Test_CommandBuilder_resolvePriorityID(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "2", "3", "4"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Priority: v,
		}
		if err := builder.ensurePriorityID(command, line); err != nil {
			t.Errorf("Could not resolve. value: %v, error: %s", line.Priority, err.Error())
		}
	}
}

func Test_CommandBuilder_resolvePriorityID_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"1", "5"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			Priority: v,
		}
		if err := builder.ensurePriorityID(command, line); err == nil {
			t.Errorf("Should return error. value: %v", line.Priority)
		}
	}
}

func Test_CommandBuilder_resolveParentIssue(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"", "EXAMPLE-1"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			ParentIssue: v,
		}
		if err := builder.ensureParentIssue(command, line); err != nil {
			t.Errorf("Could not resolve. %s", err.Error())
		}
	}
}

func Test_CommandBuilder_resolveParentIssue_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	tests := []string{"xxx"}
	for _, v := range tests {
		command := &Command{}
		line := &Line{
			ParentIssue: v,
		}
		if err := builder.ensureParentIssue(command, line); err == nil {
			t.Errorf("Should return error. value: %v", line.ParentIssue)
		}
	}
}

func Test_CommandBuilder_ensureCustomFields_should_return_error(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	command := &Command{CustomFields: make(map[CustomFieldID]interface{})}
	line := &Line{CustomFields: make(map[string]string)}
	line.CustomFields["Text"] = "apple"
	line.CustomFields["Sentence"] = "orange"
	line.CustomFields["Number"] = "abc"
	line.CustomFields["Date"] = "efg"
	line.CustomFields["Single List"] = "select-11"
	line.CustomFields["Multiple List"] = "select-12"
	line.CustomFields["Checkbox"] = "select-13"
	line.CustomFields["Radio"] = "select-13"

	if err := builder.ensureCustomFields(command, line); err == nil {
		t.Error("Should return error.")
	}
}

func Test_CommandBuilder_ensureCustomFields(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	builder := injector.Get(new(CommandBuilder)).(*commandBuilder)
	command := &Command{CustomFields: make(map[CustomFieldID]interface{})}
	line := &Line{CustomFields: make(map[string]string)}
	line.CustomFields["Text"] = "apple"
	line.CustomFields["Sentence"] = "orange"
	line.CustomFields["Number"] = "1"
	line.CustomFields["Date"] = "2018-01-01"
	line.CustomFields["Single List"] = "select-1"
	line.CustomFields["Multiple List"] = "select-2"
	line.CustomFields["Checkbox"] = "select-3"
	line.CustomFields["Radio"] = "select-3"

	if err := builder.ensureCustomFields(command, line); err != nil {
		t.Errorf("Could not resolve. %s", err.Error())
	}
	if command.CustomFields[CustomFieldID(1)] != line.CustomFields["Text"] {
		t.Errorf("Could not match. want: %v, result: %v", line.CustomFields["Text"], command.CustomFields[CustomFieldID(1)])
	}
	if command.CustomFields[CustomFieldID(2)] != line.CustomFields["Sentence"] {
		t.Errorf("Could not match. want: %v, result: %v", line.CustomFields["Sentence"], command.CustomFields[CustomFieldID(2)])
	}
	if command.CustomFields[CustomFieldID(3)] != line.CustomFields["Number"] {
		t.Errorf("Could not match. want: %v, result: %v", line.CustomFields["Number"], command.CustomFields[CustomFieldID(3)])
	}
	if command.CustomFields[CustomFieldID(4)] != line.CustomFields["Date"] {
		t.Errorf("Could not match. want: %v, result: %v", line.CustomFields["Date"], command.CustomFields[CustomFieldID(4)])
	}
	if command.CustomFields[CustomFieldID(5)] != 1 {
		t.Errorf("Could not match. want: %v, result: %v", 1, command.CustomFields[CustomFieldID(5)])
	}
	if command.CustomFields[CustomFieldID(6)] != 2 {
		t.Errorf("Could not match. want: %v, result: %v", 2, command.CustomFields[CustomFieldID(6)])
	}
	if command.CustomFields[CustomFieldID(7)] != 3 {
		t.Errorf("Could not match. want: %v, result: %v", 3, command.CustomFields[CustomFieldID(7)])
	}
	if command.CustomFields[CustomFieldID(8)] != 3 {
		t.Errorf("Could not match. want: %v, result: %v", 3, command.CustomFields[CustomFieldID(8)])
	}
}
