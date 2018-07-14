package bbir

import "strings"



func NewLine(record []string) *Line {
	line := &Line{}
	for i, v := range record {
		injectors[i](v, line)
	}
	return line
}

var injectors = []func(string, *Line){
	func(v string, line *Line) {
		line.Summary = v
	},
	func(v string, line *Line) {
		line.Description = v
	},
	func(v string, line *Line) {
		line.StartDate = strings.Replace(v, "/", "-", 3)
	},
	func(v string, line *Line) {
		line.DueDate = strings.Replace(v, "/", "-", 3)
	},
	func(v string, line *Line) {
		line.EstimatedHours = v
	},
	func(v string, line *Line) {
		line.ActualHours = v
	},
	func(v string, line *Line) {
		line.IssueType = v
	},
	func(v string, line *Line) {
		line.Category = v
	},
	func(v string, line *Line) {
		line.Version = v
	},
	func(v string, line *Line) {
		line.Milestone = v
	},
	func(v string, line *Line) {
		line.Priority = v
	},
	func(v string, line *Line) {
		line.Assignee = v
	},
	func(v string, line *Line) {
		line.ParentIssue = v
	},
}

type Line struct {
	Summary        string
	Description    string
	StartDate      string
	DueDate        string
	EstimatedHours string
	ActualHours    string
	IssueType      string
	Category       string
	Version        string
	Milestone      string
	Priority       string
	Assignee       string
	ParentIssue    string
}
