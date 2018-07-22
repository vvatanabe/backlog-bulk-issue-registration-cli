package bbir

import "fmt"

type Messages interface {
	CanNotOpenFile(path string) string
	CanNotReadFile(path string, err error) string
	SummaryIsRequired() string
	IssueTypeIsRequired() string
	IssueTypeIsNotRegistered(name string) string
	PriorityIsRequired() string
	PriorityIsInvalid(id string) string
	StartDateIsInvalid(date string) string
	DueDateIsInvalid(date string) string
	StartDateIsAfterDueDate(start, due string) string
	EstimatedHoursIsInvalid(hours string) string
	ActualHoursHoursIsInvalid(hours string) string
	CategoryIsNotRegistered(name string) string
	VersionIsNotRegistered(name string) string
	MilestoneIsNotRegistered(name string) string
	AssigneeIsNotJoining(name string) string
	ParentIssueIsNotRegistered(issueKey string) string
	ParentIssueAlreadyRegisteredAsChildIssue(issueIDOrKey interface{}) string

	FailedToCallBacklogAPI(error) string

	Line(int) string
	TagOfValidationProgressBar() string
	TagOfRegistrationProgressBar() string

	CustomFieldIsNotRegistered(name string) string
	CustomFieldValueShouldBeTypeInt(name, value string) string
	CustomFieldValueShouldBeTypeDate(name, value string) string
	CustomFieldChoiceIsNotRegistered(name string, choice string) string
}

func NewJapanese() Messages {
	return &Japanese{}
}

type Japanese struct {
}

func (m *Japanese) CanNotOpenFile(path string) string {
	return fmt.Sprintf("ファイル (%v) を開くことができません", path)
}

func (m *Japanese) CanNotReadFile(path string, err error) string {
	return fmt.Sprintf("ファイル (%v) を読み込むことができません - %v", path, err.Error())
}

func (m *Japanese) SummaryIsRequired() string {
	return "件名は必須です"
}

func (m *Japanese) IssueTypeIsRequired() string {
	return "課題の種別は必須です"
}

func (m *Japanese) IssueTypeIsNotRegistered(name string) string {
	return fmt.Sprintf("課題の種別 (%v) は登録されていません", name)
}

func (m *Japanese) PriorityIsRequired() string {
	return "課題の優先度は必須です"
}

func (m *Japanese) PriorityIsInvalid(id string) string {
	return fmt.Sprintf("課題の優先度 (%v) は不正です", id)
}

func (m *Japanese) StartDateIsInvalid(date string) string {
	return fmt.Sprintf("課題の開始日 (%v) は不正です", date)
}

func (m *Japanese) DueDateIsInvalid(date string) string {
	return fmt.Sprintf("課題の期限日 (%v) は不正です", date)
}

func (m *Japanese) StartDateIsAfterDueDate(start, due string) string {
	return fmt.Sprintf(" 課題の開始日 (%v) は期限日 (%v) より前に設定してください", start, due)
}

func (m *Japanese) EstimatedHoursIsInvalid(hours string) string {
	return fmt.Sprintf("課題の予定時間 (%v) は不正です 例. 0.25, 1.5, 1", hours)
}

func (m *Japanese) ActualHoursHoursIsInvalid(hours string) string {
	return fmt.Sprintf("課題の実績時間 (%v) は不正です 例. 0.25, 1.5, 1", hours)
}

func (m *Japanese) CategoryIsNotRegistered(name string) string {
	return fmt.Sprintf("課題のカテゴリ (%v) は登録されていません", name)
}

func (m *Japanese) VersionIsNotRegistered(name string) string {
	return fmt.Sprintf("課題のバージョン (%v) は登録されていません", name)
}

func (m *Japanese) MilestoneIsNotRegistered(name string) string {
	return fmt.Sprintf("課題のマイルストーン (%v) は登録されていません", name)
}

func (m *Japanese) AssigneeIsNotJoining(name string) string {
	return fmt.Sprintf("課題の担当者 (%v) はプロジェクトに参加していません", name)
}

func (m *Japanese) ParentIssueIsNotRegistered(issueKey string) string {
	return fmt.Sprintf("親課題 (%v) は登録されていません", issueKey)
}

func (m *Japanese) ParentIssueAlreadyRegisteredAsChildIssue(issueIDOrKey interface{}) string {
	// TODO Message is difficult to understand
	return fmt.Sprintf("親課題 (%v) は既に小課題として登録されています", issueIDOrKey)
}

func (m *Japanese) FailedToCallBacklogAPI(err error) string {
	return fmt.Sprintf("Backlog API の実行に呼び出しに失敗しました - %v", err.Error())
}

func (m *Japanese) Line(lineNumber int) string {
	return fmt.Sprintf("%v 行目", lineNumber)
}

func (m *Japanese) TagOfValidationProgressBar() string {
	return "課題データの確認"
}

func (m *Japanese) TagOfRegistrationProgressBar() string {
	return "課題データを登録"
}

func (m *Japanese) CustomFieldIsNotRegistered(name string) string {
	return fmt.Sprintf("カスタム属性 (%v) は登録されていません", name)
}

func (m *Japanese) CustomFieldValueShouldBeTypeInt(name, value string) string {
	return fmt.Sprintf("このカスタム属性 (%v) の値はint型でなければなりません - 値: %v", name, value)
}

func (m *Japanese) CustomFieldValueShouldBeTypeDate(name, value string) string {
	return fmt.Sprintf("このカスタム属性 (%v) の値はDate型でなければなりません - 値: %v", name, value)
}

func (m *Japanese) CustomFieldChoiceIsNotRegistered(name string, choice string) string {
	return fmt.Sprintf("このカスタム属性 (%v) の選択肢 (%v) は登録されていません", name, choice)
}

func NewEnglish() Messages {
	return &English{}
}

type English struct {
}

func (m *English) CanNotOpenFile(path string) string {
	return fmt.Sprintf("Can not open file (%v)", path)
}

func (m *English) CanNotReadFile(path string, err error) string {
	return fmt.Sprintf("Can not read file (%v) - %v", path, err.Error())
}

func (m *English) SummaryIsRequired() string {
	return "The summary is required"
}

func (m *English) IssueTypeIsRequired() string {
	return "The issue type is required"
}

func (m *English) IssueTypeIsNotRegistered(name string) string {
	return fmt.Sprintf("The issue type (%v) is not registered", name)
}

func (m *English) PriorityIsRequired() string {
	return "The priority is required"
}

func (m *English) PriorityIsInvalid(id string) string {
	return fmt.Sprintf("The priority (%v) is invalid", id)
}

func (m *English) StartDateIsInvalid(date string) string {
	return fmt.Sprintf("The start date (%v) is invalid", date)
}

func (m *English) DueDateIsInvalid(date string) string {
	return fmt.Sprintf("The due date (%v) is invalid", date)
}

func (m *English) StartDateIsAfterDueDate(start, due string) string {
	return fmt.Sprintf("Start Date (%v) must be set Before due date (%v)", start, due)
}

func (m *English) EstimatedHoursIsInvalid(hours string) string {
	return fmt.Sprintf("The estimated hours (%v) is invalid Ex. 0.25, 1.5, 1", hours)
}

func (m *English) ActualHoursHoursIsInvalid(hours string) string {
	return fmt.Sprintf("The actual hours (%v) is invalid Ex. 0.25, 1.5, 1", hours)
}

func (m *English) CategoryIsNotRegistered(name string) string {
	return fmt.Sprintf("The category (%v) is not registered", name)
}

func (m *English) VersionIsNotRegistered(name string) string {
	return fmt.Sprintf("The version (%v) is not registered", name)
}

func (m *English) MilestoneIsNotRegistered(name string) string {
	return fmt.Sprintf("The milestone (%v) is not registered", name)
}

func (m *English) AssigneeIsNotJoining(name string) string {
	return fmt.Sprintf("The assignee (%v) is not joining in the project", name)
}

func (m *English) ParentIssueIsNotRegistered(issueKey string) string {
	return fmt.Sprintf("The parent issue (%v) is not registered", issueKey)
}

func (m *English) ParentIssueAlreadyRegisteredAsChildIssue(issueIDOrKey interface{}) string {
	// TODO Message is difficult to understand
	return fmt.Sprintf("The parent issue (%v) has already been registered as a child issue", issueIDOrKey)
}

func (m *English) FailedToCallBacklogAPI(err error) string {
	return fmt.Sprintf("Failed to call Backlog API - %v", err.Error())
}

func (m *English) Line(lineNumber int) string {
	return fmt.Sprintf("Line %v", lineNumber)
}

func (m *English) TagOfValidationProgressBar() string {
	return "Validation   Issue"
}

func (m *English) TagOfRegistrationProgressBar() string {
	return "Registration Issue"
}

func (m *English) CustomFieldIsNotRegistered(name string) string {
	return fmt.Sprintf("The custom field (%v) is not registered", name)
}

func (m *English) CustomFieldValueShouldBeTypeInt(name, value string) string {
	return fmt.Sprintf("The value of this custom field (%v) should be type int - value: %v", name, value)
}

func (m *English) CustomFieldValueShouldBeTypeDate(name, value string) string {
	return fmt.Sprintf("The value of this custom field (%v) should be type date - value: %v", name, value)
}

func (m *English) CustomFieldChoiceIsNotRegistered(name string, choice string) string {
	return fmt.Sprintf("The choice of this custom field (%v) is not registered - choice: %v", name, choice)
}
